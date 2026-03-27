package counter

import (
	"context"
	"errors"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

const (
	keyDoesNotExist    = 20
	preconditionFailed = 22
	counterKeyPrefix   = "g-counter/"
)

type kvStore interface {
	ReadInt(ctx context.Context, key string) (int, error)
	CompareAndSwap(ctx context.Context, key string, from, to any, createIfNotExists bool) error
}

type Service struct {
	kv      kvStore
	selfID  maelstromx.NodeIDFunc
	nodeIDs maelstromx.NodeIDsFunc
}

func NewService(kv kvStore, selfID maelstromx.NodeIDFunc, nodeIDs maelstromx.NodeIDsFunc) *Service {
	return &Service{
		kv:      kv,
		selfID:  selfID,
		nodeIDs: nodeIDs,
	}
}

func (s *Service) Add(ctx context.Context, delta int) error {
	key := counterKey(s.selfID())

	for {
		current, err := s.kv.ReadInt(ctx, key)
		if err != nil {
			if isRPCCode(err, keyDoesNotExist) {
				if err := s.kv.CompareAndSwap(ctx, key, 0, delta, true); err != nil {
					if isRetryableCAS(err) {
						continue
					}

					return err
				}

				return nil
			}

			return err
		}

		if err := s.kv.CompareAndSwap(ctx, key, current, current+delta, false); err != nil {
			if isRetryableCAS(err) {
				continue
			}

			return err
		}

		return nil
	}
}

func (s *Service) Read(ctx context.Context) (int, error) {
	total := 0
	for _, nodeID := range s.nodeIDs() {
		value, err := s.kv.ReadInt(ctx, counterKey(nodeID))
		if err != nil {
			if isRPCCode(err, keyDoesNotExist) {
				continue
			}

			return 0, err
		}

		total += value
	}

	return total, nil
}

func counterKey(nodeID string) string {
	return fmt.Sprintf("%s%s", counterKeyPrefix, nodeID)
}

func isRetryableCAS(err error) bool {
	return isRPCCode(err, keyDoesNotExist) || isRPCCode(err, preconditionFailed)
}

func isRPCCode(err error, code int) bool {
	var rpcErr *maelstrom.RPCError
	return errors.As(err, &rpcErr) && rpcErr.Code == code
}
