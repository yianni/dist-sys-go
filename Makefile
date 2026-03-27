BIN_DIR := bin
MAELSTROM_DIR := tools/maelstrom
MAELSTROM := $(MAELSTROM_DIR)/maelstrom

.PHONY: fmt test race build-echo build-uniqueids build-broadcast build-counter build-kafka build-txn build setup-maelstrom maelstrom-echo maelstrom-uniqueids maelstrom-broadcast maelstrom-counter maelstrom-kafka maelstrom-txn validate

fmt:
	gofmt -w cmd internal

test:
	go test ./...

race:
	go test -race ./...

build-echo:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/echo ./cmd/echo

build-uniqueids:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/uniqueids ./cmd/uniqueids

build-broadcast:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/broadcast ./cmd/broadcast

build-counter:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/counter ./cmd/counter

build-kafka:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/kafka ./cmd/kafka

build-txn:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/txn ./cmd/txn

build: build-echo build-uniqueids build-broadcast build-counter build-kafka build-txn

setup-maelstrom:
	mkdir -p tools
	curl -L -o tools/maelstrom.tar.bz2 https://github.com/jepsen-io/maelstrom/releases/download/v0.2.4/maelstrom.tar.bz2
	tar -xjf tools/maelstrom.tar.bz2 -C tools

maelstrom-echo: build-echo
	$(MAELSTROM) test -w echo --bin ./$(BIN_DIR)/echo --node-count 1 --time-limit 10

maelstrom-uniqueids: build-uniqueids
	$(MAELSTROM) test -w unique-ids --bin ./$(BIN_DIR)/uniqueids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

maelstrom-broadcast: build-broadcast
	$(MAELSTROM) test -w broadcast --bin ./$(BIN_DIR)/broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

maelstrom-counter: build-counter
	$(MAELSTROM) test -w g-counter --bin ./$(BIN_DIR)/counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition

maelstrom-kafka: build-kafka
	$(MAELSTROM) test -w kafka --bin ./$(BIN_DIR)/kafka --node-count 2 --concurrency 2n --time-limit 20 --rate 1000

maelstrom-txn: build-txn
	$(MAELSTROM) test -w txn-rw-register --bin ./$(BIN_DIR)/txn --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition

validate: maelstrom-echo maelstrom-uniqueids maelstrom-broadcast maelstrom-counter maelstrom-kafka maelstrom-txn
