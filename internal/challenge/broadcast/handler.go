package broadcast

import (
	"context"
	"slices"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

const maxGossipBatchSize = 128
const fullSyncEvery = 20

type Handler struct {
	node    *maelstrom.Node
	service *Service
	ticks   int
}

func NewHandler(node *maelstrom.Node, service *Service) *Handler {
	return &Handler{
		node:    node,
		service: service,
	}
}

func (h *Handler) Start(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.gossipTick()
			}
		}
	}()
}

func (h *Handler) HandleBroadcast(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[broadcastRequest](msg)
	if err != nil {
		return err
	}

	firstSeen := h.service.Add(req.Message)
	if firstSeen {
		go h.sendGossip([]int{req.Message})
	}

	return h.node.Reply(msg, broadcastResponse{Type: "broadcast_ok"})
}

func (h *Handler) HandleRead(msg maelstrom.Message) error {
	_, err := maelstromx.Decode[readRequest](msg)
	if err != nil {
		return err
	}

	return h.node.Reply(msg, readResponse{
		Type:     "read_ok",
		Messages: h.service.Messages(),
	})
}

func (h *Handler) HandleTopology(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[topologyRequest](msg)
	if err != nil {
		return err
	}

	h.service.ConfigureTopology(req.Topology)

	return h.node.Reply(msg, topologyResponse{Type: "topology_ok"})
}

func (h *Handler) HandleGossip(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[gossipRequest](msg)
	if err != nil {
		return err
	}

	if newMessages := h.service.Merge(req.Messages); len(newMessages) > 0 {
		go h.sendGossip(newMessages)
	}

	return nil
}

func (h *Handler) gossip() {
	messages := h.service.Messages()
	if len(messages) == 0 {
		return
	}

	h.sendGossip(messages)
}

func (h *Handler) gossipTick() {
	h.ticks++
	if h.ticks%fullSyncEvery == 0 {
		h.gossip()
		return
	}

	if dirty := h.service.DrainDirty(); len(dirty) > 0 {
		h.sendGossip(dirty)
	}
}

func (h *Handler) sendGossip(messages []int) {
	for start := 0; start < len(messages); start += maxGossipBatchSize {
		end := min(start+maxGossipBatchSize, len(messages))
		batch := slices.Clone(messages[start:end])
		body := gossipRequest{
			Type:     "gossip",
			Messages: batch,
		}

		for _, peer := range h.service.Peers() {
			_ = h.node.Send(peer, body)
		}
	}
}
