package txn

import (
	"context"
	"slices"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

const maxSyncBatchSize = 128
const fullSyncEvery = 20

type Handler struct {
	node    *maelstrom.Node
	service *Service
	ticks   int
}

func NewHandler(node *maelstrom.Node, service *Service) *Handler {
	return &Handler{node: node, service: service}
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

func (h *Handler) HandleTxn(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[txnRequest](msg)
	if err != nil {
		return err
	}

	result, writes := h.service.Apply(req.Txn)
	if len(writes) > 0 {
		go h.gossipWrites(writes)
	}

	return h.node.Reply(msg, txnResponse{Type: "txn_ok", Txn: result})
}

func (h *Handler) HandleTxnSync(msg maelstrom.Message) error {
	req, err := maelstromx.Decode[syncRequest](msg)
	if err != nil {
		return err
	}

	h.service.Merge(req.Writes)
	return nil
}

func (h *Handler) gossipAll() {
	writes := h.service.SnapshotWrites()
	if len(writes) == 0 {
		return
	}
	h.gossipWrites(writes)
}

func (h *Handler) gossipTick() {
	h.ticks++
	if h.ticks%fullSyncEvery == 0 {
		h.gossipAll()
		return
	}

	if writes := h.service.DrainDirtyWrites(); len(writes) > 0 {
		h.gossipWrites(writes)
	}
}

func (h *Handler) gossipWrites(writes []writeState) {
	for start := 0; start < len(writes); start += maxSyncBatchSize {
		end := min(start+maxSyncBatchSize, len(writes))
		batch := slices.Clone(writes[start:end])
		body := syncRequest{Type: "txn_sync", Writes: batch}
		for _, peer := range h.service.Peers() {
			_ = h.node.Send(peer, body)
		}
	}
}
