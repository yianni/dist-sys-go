package main

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/challenge/txn"
	"dist-sys-go/internal/platform/maelstromx"
)

func main() {
	node := maelstrom.NewNode()

	service := txn.NewService(node.ID, node.NodeIDs)
	handler := txn.NewHandler(node, service)
	handler.Start(context.Background(), 250*time.Millisecond)

	node.Handle("txn", handler.HandleTxn)
	node.Handle("txn_sync", handler.HandleTxnSync)

	maelstromx.MustRun(node)
}
