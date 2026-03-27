package uniqueids

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

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

func (h Handler) HandleGenerate(msg maelstrom.Message) error {
	_, err := maelstromx.Decode[generateRequest](msg)
	if err != nil {
		return err
	}

	resp := generateResponse{
		Type: "generate_ok",
		ID:   h.service.Generate(),
	}

	return h.node.Reply(msg, resp)
}
