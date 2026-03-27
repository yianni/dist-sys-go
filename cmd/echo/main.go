package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/challenge/echo"
	"dist-sys-go/internal/platform/maelstromx"
)

func main() {
	node := maelstrom.NewNode()
	service := echo.NewService()
	handler := echo.NewHandler(node, service)

	node.Handle("echo", handler.HandleEcho)

	maelstromx.MustRun(node)
}
