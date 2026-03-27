package broadcast

import (
	"sort"
	"sync"

	"dist-sys-go/internal/platform/maelstromx"
)

type Service struct {
	selfID  maelstromx.NodeIDFunc
	nodeIDs maelstromx.NodeIDsFunc

	mu       sync.RWMutex
	messages map[int]struct{}
	dirty    map[int]struct{}
	peers    []string
}

func NewService(selfID maelstromx.NodeIDFunc, nodeIDs maelstromx.NodeIDsFunc) *Service {
	return &Service{
		selfID:   selfID,
		nodeIDs:  nodeIDs,
		messages: make(map[int]struct{}),
		dirty:    make(map[int]struct{}),
	}
}

func (s *Service) Add(message int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.messages[message]; exists {
		return false
	}

	s.messages[message] = struct{}{}
	s.dirty[message] = struct{}{}
	return true
}

func (s *Service) Merge(messages []int) []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	newMessages := make([]int, 0, len(messages))
	for _, message := range messages {
		if _, exists := s.messages[message]; exists {
			continue
		}

		s.messages[message] = struct{}{}
		s.dirty[message] = struct{}{}
		newMessages = append(newMessages, message)
	}

	sort.Ints(newMessages)
	return newMessages
}

func (s *Service) Messages() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]int, 0, len(s.messages))
	for message := range s.messages {
		result = append(result, message)
	}

	sort.Ints(result)
	return result
}

func (s *Service) DrainDirty() []int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.dirty) == 0 {
		return nil
	}

	result := make([]int, 0, len(s.dirty))
	for message := range s.dirty {
		result = append(result, message)
	}

	clear(s.dirty)
	sort.Ints(result)
	return result
}

func (s *Service) ConfigureTopology(topology map[string][]string) {
	self := s.selfID()
	peers := uniqueStrings(topology[self])
	if len(peers) == 0 {
		peers = s.defaultPeers()
	}

	s.mu.Lock()
	s.peers = peers
	s.mu.Unlock()
}

func (s *Service) Peers() []string {
	s.mu.RLock()
	if len(s.peers) > 0 {
		peers := append([]string(nil), s.peers...)
		s.mu.RUnlock()
		return peers
	}
	s.mu.RUnlock()

	return s.defaultPeers()
}

func (s *Service) defaultPeers() []string {
	return uniqueStrings(maelstromx.Peers(s.selfID, s.nodeIDs))
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	sort.Strings(result)
	return result
}
