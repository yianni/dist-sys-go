---- MODULE Broadcast ----
EXTENDS FiniteSets, TLC

CONSTANT Nodes
CONSTANT Messages

Packet == [from : Nodes, to : Nodes, msgs : SUBSET Messages]

VARIABLE known
VARIABLE inflight
VARIABLE everKnown
VARIABLE broadcasted

Init ==
  /\ known = [n \in Nodes |-> {}]
  /\ inflight = {}
  /\ everKnown = known
  /\ broadcasted = {}

Broadcast(n, m) ==
  /\ n \in Nodes
  /\ m \in Messages
  /\ known' = [known EXCEPT ![n] = @ \cup {m}]
  /\ inflight' = inflight
  /\ everKnown' = [node \in Nodes |-> everKnown[node] \cup known'[node]]
  /\ broadcasted' = broadcasted \cup {m}

SendGossip(from, to) ==
  /\ from \in Nodes
  /\ to \in Nodes
  /\ from # to
  /\ inflight' = inflight \cup {[from |-> from, to |-> to, msgs |-> known[from]]}
  /\ known' = known
  /\ everKnown' = [node \in Nodes |-> everKnown[node] \cup known'[node]]
  /\ broadcasted' = broadcasted

RecvGossip(pkt) ==
  /\ pkt \in inflight
  /\ known' = [known EXCEPT ![pkt.to] = @ \cup pkt.msgs]
  /\ inflight' = inflight \ {pkt}
  /\ everKnown' = [node \in Nodes |-> everKnown[node] \cup known'[node]]
  /\ broadcasted' = broadcasted

Next ==
  \/ \E n \in Nodes:
       \E m \in Messages:
         Broadcast(n, m)
  \/ \E from \in Nodes:
       \E to \in Nodes:
         SendGossip(from, to)
  \/ \E pkt \in inflight:
       RecvGossip(pkt)

TypeInvariant ==
  /\ known \in [Nodes -> SUBSET Messages]
  /\ inflight \subseteq Packet
  /\ everKnown \in [Nodes -> SUBSET Messages]
  /\ broadcasted \subseteq Messages

NoForgetting ==
  \A n \in Nodes:
    everKnown[n] \subseteq known[n]

AllKnown ==
  UNION {known[n] : n \in Nodes}

NoGhostKnowledge ==
  AllKnown \subseteq broadcasted

====
