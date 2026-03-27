package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/challenge/uniqueids"
	"dist-sys-go/internal/platform/maelstromx"
)

func main() {
	node := maelstrom.NewNode()
	service := uniqueids.NewService(node.ID)
	handler := uniqueids.NewHandler(node, service)

	node.Handle("generate", handler.HandleGenerate)

	maelstromx.MustRun(node)
}
