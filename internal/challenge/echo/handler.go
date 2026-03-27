package echo

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

func (h Handler) HandleEcho(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[echoRequest](msg)
	if err != nil {
		return err
	}

	resp := echoResponse{
		Type: "echo_ok",
		Echo: h.service.Echo(req.Echo),
	}

	return h.node.Reply(msg, resp)
}
