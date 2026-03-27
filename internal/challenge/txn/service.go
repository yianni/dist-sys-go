package txn

import (
	"sort"
	"sync"
	"sync/atomic"

	"dist-sys-go/internal/platform/maelstromx"
)

type version struct {
	Counter uint64 `json:"counter"`
	NodeID  string `json:"node_id"`
}

type registerState struct {
	Value   int     `json:"value"`
	Version version `json:"version"`
}

type writeState struct {
	Key     int     `json:"key"`
	Value   int     `json:"value"`
	Version version `json:"version"`
}

type Service struct {
	selfID  maelstromx.NodeIDFunc
	nodeIDs maelstromx.NodeIDsFunc
	clock   atomic.Uint64

	mu    sync.RWMutex
	store map[int]registerState
	peers []string
}

func NewService(selfID maelstromx.NodeIDFunc, nodeIDs maelstromx.NodeIDsFunc) *Service {
	return &Service{
		selfID:  selfID,
		nodeIDs: nodeIDs,
		store:   make(map[int]registerState),
	}
}

func (s *Service) Apply(txn []operation) ([]operation, []writeState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]operation, 0, len(txn))
	writes := make([]writeState, 0)

	for _, op := range txn {
		switch op.Kind {
		case "r":
			read := operation{Kind: op.Kind, Key: op.Key}
			if state, ok := s.store[op.Key]; ok {
				value := state.Value
				read.Value = &value
			}
			result = append(result, read)
		case "w":
			if op.Value == nil {
				result = append(result, operation{Kind: op.Kind, Key: op.Key})
				continue
			}

			ver := s.nextVersion()
			state := registerState{Value: *op.Value, Version: ver}
			s.store[op.Key] = state
			value := *op.Value
			result = append(result, operation{Kind: op.Kind, Key: op.Key, Value: &value})
			writes = append(writes, writeState{Key: op.Key, Value: value, Version: ver})
		}
	}

	return result, writes
}

func (s *Service) Merge(writes []writeState) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	changed := false
	for _, write := range writes {
		s.observeVersion(write.Version)

		current, ok := s.store[write.Key]
		if ok && !write.Version.After(current.Version) {
			continue
		}

		s.store[write.Key] = registerState{Value: write.Value, Version: write.Version}
		changed = true
	}

	return changed
}

func (s *Service) SnapshotWrites() []writeState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]int, 0, len(s.store))
	for key := range s.store {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	writes := make([]writeState, 0, len(keys))
	for _, key := range keys {
		state := s.store[key]
		writes = append(writes, writeState{Key: key, Value: state.Value, Version: state.Version})
	}

	return writes
}

func (s *Service) Peers() []string {
	s.mu.RLock()
	if len(s.peers) > 0 {
		peers := append([]string(nil), s.peers...)
		s.mu.RUnlock()
		return peers
	}
	s.mu.RUnlock()

	peers := maelstromx.Peers(s.selfID, s.nodeIDs)

	s.mu.Lock()
	s.peers = append([]string(nil), peers...)
	s.mu.Unlock()

	return peers
}

func (v version) After(other version) bool {
	if v.Counter != other.Counter {
		return v.Counter > other.Counter
	}
	return v.NodeID > other.NodeID
}

func (s *Service) nextVersion() version {
	return version{
		Counter: s.clock.Add(1),
		NodeID:  s.selfID(),
	}
}

func (s *Service) observeVersion(v version) {
	for {
		current := s.clock.Load()
		if current >= v.Counter {
			return
		}

		if s.clock.CompareAndSwap(current, v.Counter) {
			return
		}
	}
}
