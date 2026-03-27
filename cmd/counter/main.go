package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/challenge/counter"
	"dist-sys-go/internal/platform/maelstromx"
)

func main() {
	node := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(node)

	service := counter.NewService(kv, node.ID, node.NodeIDs)
	handler := counter.NewHandler(node, service)

	node.Handle("add", handler.HandleAdd)
	node.Handle("read", handler.HandleRead)

	maelstromx.MustRun(node)
}
