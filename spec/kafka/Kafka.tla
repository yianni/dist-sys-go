---- MODULE Kafka ----
EXTENDS Sequences, Naturals, Integers, TLC

CONSTANT Nodes
CONSTANT Keys
CONSTANT Msgs
CONSTANT MaxLogLen

NoOffset == -1

VARIABLE logs
VARIABLE committed
VARIABLE previousLogs
VARIABLE previousCommitted

Init ==
  /\ logs = [k \in Keys |-> <<>>]
  /\ committed = [k \in Keys |-> NoOffset]
  /\ previousLogs = logs
  /\ previousCommitted = committed

AppendMessage(k, msg) ==
  /\ k \in Keys
  /\ msg \in Msgs
  /\ Len(logs[k]) < MaxLogLen
  /\ logs' = [logs EXCEPT ![k] = Append(@, msg)]
  /\ committed' = committed
  /\ previousLogs' = logs
  /\ previousCommitted' = committed

Commit(k, off) ==
  /\ k \in Keys
  /\ Len(logs[k]) > 0
  /\ off \in 0..(Len(logs[k]) - 1)
  /\ IF committed[k] = NoOffset THEN TRUE ELSE off >= committed[k]
  /\ committed' = [committed EXCEPT ![k] = off]
  /\ logs' = logs
  /\ previousLogs' = logs
  /\ previousCommitted' = committed

Next ==
  \/ \E k \in Keys:
       \E msg \in Msgs:
         AppendMessage(k, msg)
  \/ \E k \in Keys:
       \E off \in 0..(IF Len(logs[k]) = 0 THEN 0 ELSE Len(logs[k]) - 1):
         Commit(k, off)

TypeInvariant ==
  /\ logs \in [Keys -> Seq(Msgs)]
  /\ committed \in [Keys -> (Nat \cup {NoOffset})]
  /\ previousLogs \in [Keys -> Seq(Msgs)]
  /\ previousCommitted \in [Keys -> (Nat \cup {NoOffset})]

CommittedOffsetValid ==
  \A k \in Keys:
    IF committed[k] = NoOffset THEN TRUE ELSE committed[k] < Len(logs[k])

PollView(k, off) ==
  IF off >= Len(logs[k]) THEN
    <<>>
  ELSE
    [i \in (off + 1)..Len(logs[k]) |-> <<i - 1, logs[k][i]>>]

IsPrefix(prefix, seq) ==
  /\ Len(prefix) <= Len(seq)
  /\ SubSeq(seq, 1, Len(prefix)) = prefix

AppendOnly ==
  \A k \in Keys:
    IsPrefix(previousLogs[k], logs[k])

CommittedMonotonic ==
  \A k \in Keys:
    IF previousCommitted[k] = NoOffset THEN TRUE ELSE committed[k] >= previousCommitted[k]

PollViewValid ==
  \A k \in Keys:
    \A off \in 0..Len(logs[k]):
      \A i \in DOMAIN PollView(k, off):
        /\ PollView(k, off)[i][1] >= off
        /\ PollView(k, off)[i][1] < Len(logs[k])
        /\ PollView(k, off)[i][2] \in Msgs

====
