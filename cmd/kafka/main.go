package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/challenge/kafka"
	"dist-sys-go/internal/platform/maelstromx"
)

func main() {
	node := maelstrom.NewNode()

	service := kafka.NewService(node.ID, node.NodeIDs)
	handler := kafka.NewHandler(node, service)

	node.Handle("send", handler.HandleSend)
	node.Handle("poll", handler.HandlePoll)
	node.Handle("commit_offsets", handler.HandleCommitOffsets)
	node.Handle("list_committed_offsets", handler.HandleListCommittedOffsets)
	node.Handle("send_local", handler.HandleSendLocal)
	node.Handle("poll_local", handler.HandlePollLocal)
	node.Handle("commit_offsets_local", handler.HandleCommitOffsetsLocal)
	node.Handle("list_committed_offsets_local", handler.HandleListCommittedOffsetsLocal)

	maelstromx.MustRun(node)
}
