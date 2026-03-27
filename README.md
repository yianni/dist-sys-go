# dist-sys-go

Go implementations of the Fly.io Gossip Glomers distributed systems challenges.

## Layout

```text
cmd/
  echo/           # challenge entrypoint
  uniqueids/      # challenge entrypoint
  broadcast/      # challenge entrypoint
  counter/        # challenge entrypoint
  kafka/          # challenge entrypoint
  txn/            # challenge entrypoint
internal/
  challenge/
    echo/         # challenge-specific logic and message handling
    uniqueids/    # challenge-specific logic and message handling
    broadcast/    # challenge-specific logic and message handling
    counter/      # challenge-specific logic and message handling
    kafka/        # challenge-specific logic and message handling
    txn/          # challenge-specific logic and message handling
  platform/
    maelstromx/   # shared Maelstrom helpers
```

## Quickstart

```bash
make fmt
make test
make build
make setup-maelstrom
make validate
```

## Challenges

| Challenge                                                    | Scope   | Approach                                              | Validate                   |
|--------------------------------------------------------------|---------|-------------------------------------------------------|----------------------------|
| [Echo](https://fly.io/dist-sys/1)                            | `1`     | request/reply baseline                                | `make maelstrom-echo`      |
| [Unique ID Generation](https://fly.io/dist-sys/2)            | `2`     | `node-id + counter`                                   | `make maelstrom-uniqueids` |
| [Broadcast](https://fly.io/dist-sys/3a)                      | `3a-3c` | gossip plus periodic anti-entropy                     | `make maelstrom-broadcast` |
| [Grow-Only Counter](https://fly.io/dist-sys/4)               | `4`     | per-node `seq-kv` counters, summed on read            | `make maelstrom-counter`   |
| [Kafka-Style Log](https://fly.io/dist-sys/5a)                | `5a-5c` | deterministic owner routing with local in-memory logs | `make maelstrom-kafka`     |
| [Totally-Available Transactions](https://fly.io/dist-sys/6a) | `6a-6b` | local txn application with async LWW replication      | `make maelstrom-txn`       |

Direct Maelstrom commands:

```bash
./tools/maelstrom/maelstrom test -w echo --bin ./bin/echo --node-count 1 --time-limit 10
./tools/maelstrom/maelstrom test -w unique-ids --bin ./bin/uniqueids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
./tools/maelstrom/maelstrom test -w broadcast --bin ./bin/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition
./tools/maelstrom/maelstrom test -w g-counter --bin ./bin/counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition
./tools/maelstrom/maelstrom test -w kafka --bin ./bin/kafka --node-count 2 --concurrency 2n --time-limit 20 --rate 1000
./tools/maelstrom/maelstrom test -w txn-rw-register --bin ./bin/txn --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition
```
