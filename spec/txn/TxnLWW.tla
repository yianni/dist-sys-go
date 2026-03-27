---- MODULE TxnLWW ----
EXTENDS Naturals, TLC

CONSTANT Nodes
CONSTANT Keys
CONSTANT Values
CONSTANT MaxClock
CONSTANT NoValue

VersionSet == [counter : Nat, node : Nodes]
NullVersion == [counter |-> 0, node |-> CHOOSE n \in Nodes: TRUE]
StoreState == [value : Values \cup {NoValue}, version : VersionSet]
ReplicationMessage == [from : Nodes, to : Nodes, key : Keys, state : StoreState]
VARIABLE store
VARIABLE clock
VARIABLE inflight

Init ==
  /\ store = [n \in Nodes |-> [k \in Keys |-> [value |-> NoValue, version |-> NullVersion]]]
  /\ clock = [n \in Nodes |-> 0]
  /\ inflight = {}

VersionAfter(v1, v2) ==
  IF v1.counter # v2.counter THEN
    v1.counter > v2.counter
  ELSE
    v1.node > v2.node

MergeState(local, remote) ==
  IF local.value = NoValue THEN
    remote
  ELSE IF VersionAfter(remote.version, local.version) THEN
    remote
  ELSE
    local

MaxCounter(a, b) ==
  IF a >= b THEN a ELSE b

ApplyWrite(n, k, v) ==
  /\ n \in Nodes
  /\ k \in Keys
  /\ v \in Values
  /\ clock[n] < MaxClock
  /\ LET nextState == [value |-> v, version |-> [counter |-> clock[n] + 1, node |-> n]]
     IN inflight' = inflight \cup { [from |-> n, to |-> peer, key |-> k, state |-> nextState] : peer \in (Nodes \ {n}) }
  /\ clock' = [clock EXCEPT ![n] = @ + 1]
  /\ store' = [store EXCEPT ![n][k] = [value |-> v, version |-> [counter |-> clock[n] + 1, node |-> n]]]

Deliver(msg) ==
  /\ msg \in inflight
  /\ store' = [store EXCEPT ![msg.to][msg.key] = MergeState(@, msg.state)]
  /\ clock' = [clock EXCEPT ![msg.to] = MaxCounter(@, msg.state.version.counter)]
  /\ inflight' = inflight \ {msg}

Next ==
  \/ \E n \in Nodes:
       \E k \in Keys:
         \E v \in Values:
           ApplyWrite(n, k, v)
  \/ \E msg \in inflight:
       Deliver(msg)

TypeInvariant ==
  /\ store \in [Nodes -> [Keys -> StoreState]]
  /\ clock \in [Nodes -> Nat]
  /\ inflight \subseteq ReplicationMessage

ObservedVersionsDoNotExceedClock ==
  \A n \in Nodes:
    \A k \in Keys:
      store[n][k].version.counter <= clock[n]

ReadView(n, k) ==
  store[n][k].value

ConvergedOnKey(k) ==
  \A n1 \in Nodes:
    \A n2 \in Nodes:
      store[n1][k] = store[n2][k]

Converged ==
  \A k \in Keys:
    ConvergedOnKey(k)

Quiescent ==
  inflight = {}

QuiescentImpliesConverged ==
  Quiescent => Converged

====
