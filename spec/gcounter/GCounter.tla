---- MODULE GCounter ----
EXTENDS Naturals, FiniteSets, TLC

CONSTANT Nodes
CONSTANT MaxDelta
CONSTANT MaxCount

VARIABLE counter
VARIABLE previousCounter

RECURSIVE SumTo(_, _)

Init ==
  /\ counter = [n \in Nodes |-> 0]
  /\ previousCounter = counter

Add(n, d) ==
  /\ n \in Nodes
  /\ d \in 1..MaxDelta
  /\ counter[n] + d <= MaxCount
  /\ counter' = [counter EXCEPT ![n] = @ + d]
  /\ previousCounter' = counter

Next ==
  \E n \in Nodes:
    \E d \in 1..MaxDelta:
      Add(n, d)

TypeInvariant ==
  /\ counter \in [Nodes -> Nat]
  /\ previousCounter \in [Nodes -> Nat]

NumericNodeDomain ==
  Nodes = 1..Cardinality(Nodes)

ComponentMonotonic ==
  \A n \in Nodes:
    counter[n] >= previousCounter[n]

SumTo(c, n) ==
  IF n = 0 THEN
    0
  ELSE
    c[n] + SumTo(c, n - 1)

Total(c) ==
  SumTo(c, Cardinality(Nodes))

ReadView ==
  Total(counter)

TotalMonotonic ==
  Total(counter) >= Total(previousCounter)

AtMostOneComponentChanges ==
  Cardinality({n \in Nodes: counter[n] # previousCounter[n]}) <= 1

ComponentValue(n) ==
  counter[n]

====
