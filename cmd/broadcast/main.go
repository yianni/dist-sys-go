package main

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/challenge/broadcast"
	"dist-sys-go/internal/platform/maelstromx"
)

func main() {
	node := maelstrom.NewNode()
	service := broadcast.NewService(node.ID, node.NodeIDs)
	handler := broadcast.NewHandler(node, service)

	handler.Start(context.Background(), 250*time.Millisecond)

	node.Handle("broadcast", handler.HandleBroadcast)
	node.Handle("read", handler.HandleRead)
	node.Handle("topology", handler.HandleTopology)
	node.Handle("gossip", handler.HandleGossip)

	maelstromx.MustRun(node)
}
