package kafka

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

type Handler struct {
	node        *maelstrom.Node
	service     *Service
	coordinator Coordinator
}

func NewHandler(node *maelstrom.Node, service *Service) Handler {
	return Handler{
		node:        node,
		service:     service,
		coordinator: NewCoordinator(node, service),
	}
}

func (h Handler) HandleSend(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[sendRequest](msg)
	if err != nil {
		return err
	}

	resp, err := h.coordinator.Send(req.Key, req.Msg)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, resp)
}

func (h Handler) HandlePoll(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[pollRequest](msg)
	if err != nil {
		return err
	}

	resp, err := h.coordinator.Poll(req.Offsets)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, resp)
}

func (h Handler) HandleCommitOffsets(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[commitOffsetsRequest](msg)
	if err != nil {
		return err
	}

	resp, err := h.coordinator.CommitOffsets(req.Offsets)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, resp)
}

func (h Handler) HandleListCommittedOffsets(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[listCommittedOffsetsRequest](msg)
	if err != nil {
		return err
	}

	resp, err := h.coordinator.ListCommittedOffsets(req.Keys)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, resp)
}

func (h Handler) HandleSendLocal(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[sendLocalRequest](msg)
	if err != nil {
		return err
	}

	offset := h.service.Send(req.Key, req.Msg)
	return h.node.Reply(msg, sendResponse{Type: "send_ok", Offset: offset})
}

func (h Handler) HandlePollLocal(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[pollLocalRequest](msg)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, pollResponse{Type: "poll_ok", Msgs: h.service.Poll(req.Offsets)})
}

func (h Handler) HandleCommitOffsetsLocal(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[commitOffsetsLocalRequest](msg)
	if err != nil {
		return err
	}

	h.service.CommitOffsets(req.Offsets)
	return h.node.Reply(msg, commitOffsetsResponse{Type: "commit_offsets_ok"})
}

func (h Handler) HandleListCommittedOffsetsLocal(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[listCommittedOffsetsLocalRequest](msg)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, listCommittedOffsetsResponse{
		Type:    "list_committed_offsets_ok",
		Offsets: h.service.ListCommittedOffsets(req.Keys),
	})
}
