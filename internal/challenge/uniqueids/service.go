package uniqueids

import (
	"fmt"
	"sync/atomic"

	"dist-sys-go/internal/platform/maelstromx"
)

type Service struct {
	nodeID  maelstromx.NodeIDFunc
	counter atomic.Uint64
}

func NewService(nodeID maelstromx.NodeIDFunc) *Service {
	return &Service{nodeID: nodeID}
}

func (s *Service) Generate() string {
	next := s.counter.Add(1)
	return fmt.Sprintf("%s-%d", s.nodeID(), next)
}
