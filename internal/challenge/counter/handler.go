package counter

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

const requestTimeout = 5 * time.Second

type Handler struct {
	node    *maelstrom.Node
	service *Service
}

func NewHandler(node *maelstrom.Node, service *Service) Handler {
	return Handler{
		node:    node,
		service: service,
	}
}

func (h Handler) HandleAdd(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[addRequest](msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if err := h.service.Add(ctx, req.Delta); err != nil {
		return err
	}

	return h.node.Reply(msg, addResponse{Type: "add_ok"})
}

func (h Handler) HandleRead(msg maelstrom.Message) error {
	_, err := maelstromx.Decode[readRequest](msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	value, err := h.service.Read(ctx)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, readResponse{
		Type:  "read_ok",
		Value: value,
	})
}
