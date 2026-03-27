package kafka

import (
	"hash/fnv"
	"sort"
	"sync"

	"dist-sys-go/internal/platform/maelstromx"
)

type streamState struct {
	records          []logRecord
	committedOffset  int
	hasCommittedRead bool
}

type Service struct {
	selfID  maelstromx.NodeIDFunc
	nodeIDs maelstromx.NodeIDsFunc

	mu      sync.RWMutex
	streams map[string]*streamState
}

func NewService(selfID maelstromx.NodeIDFunc, nodeIDs maelstromx.NodeIDsFunc) *Service {
	return &Service{
		selfID:  selfID,
		nodeIDs: nodeIDs,
		streams: make(map[string]*streamState),
	}
}

func (s *Service) Owner(key string) string {
	nodes := maelstromx.SortedNodeIDs(s.nodeIDs)
	if len(nodes) == 0 {
		return ""
	}

	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(key))
	index := int(hasher.Sum32() % uint32(len(nodes)))

	return nodes[index]
}

func (s *Service) IsOwner(key string) bool {
	return s.Owner(key) == s.selfID()
}

func (s *Service) Send(key string, msg int) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	stream := s.stream(key)
	offset := len(stream.records)
	stream.records = append(stream.records, logRecord{offset, msg})

	return offset
}

func (s *Service) Poll(offsets map[string]int) map[string][]logRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := sortedMapKeys(offsets)
	msgs := make(map[string][]logRecord, len(keys))
	for _, key := range keys {
		stream, ok := s.streams[key]
		if !ok {
			continue
		}

		start := offsets[key]
		if start < 0 || start >= len(stream.records) {
			continue
		}

		records := append([]logRecord(nil), stream.records[start:]...)
		if len(records) > 0 {
			msgs[key] = records
		}
	}

	return msgs
}

func (s *Service) CommitOffsets(offsets map[string]int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range sortedMapKeys(offsets) {
		stream := s.stream(key)
		offset := offsets[key]
		if !stream.hasCommittedRead || offset > stream.committedOffset {
			stream.committedOffset = offset
			stream.hasCommittedRead = true
		}
	}
}

func (s *Service) ListCommittedOffsets(keys []string) map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ordered := append([]string(nil), keys...)
	sort.Strings(ordered)

	offsets := make(map[string]int, len(ordered))
	for _, key := range ordered {
		stream, ok := s.streams[key]
		if !ok || !stream.hasCommittedRead {
			continue
		}

		offsets[key] = stream.committedOffset
	}

	return offsets
}

func (s *Service) stream(key string) *streamState {
	stream, ok := s.streams[key]
	if ok {
		return stream
	}

	stream = &streamState{}
	s.streams[key] = stream
	return stream
}

func sortedMapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
