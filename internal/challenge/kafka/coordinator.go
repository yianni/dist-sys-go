package kafka

import (
	"context"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"dist-sys-go/internal/platform/maelstromx"
)

const requestTimeout = 5 * time.Second

type Coordinator struct {
	node    *maelstrom.Node
	service *Service
}

func NewCoordinator(node *maelstrom.Node, service *Service) Coordinator {
	return Coordinator{node: node, service: service}
}

func (c Coordinator) Send(key string, msg int) (sendResponse, error) {
	if c.service.IsOwner(key) {
		offset := c.service.Send(key, msg)
		return sendResponse{Type: "send_ok", Offset: offset}, nil
	}

	return proxy[sendResponse](c.node, c.service.Owner(key), sendLocalRequest{
		Type: "send_local",
		Key:  key,
		Msg:  msg,
	})
}

func (c Coordinator) Poll(offsets map[string]int) (pollResponse, error) {
	groups := c.groupOffsetsByOwner(offsets)
	msgs := make(map[string][]logRecord, len(offsets))
	for owner, ownedOffsets := range groups {
		var partial map[string][]logRecord
		if owner == c.node.ID() {
			partial = c.service.Poll(ownedOffsets)
		} else {
			resp, err := proxy[pollResponse](c.node, owner, pollLocalRequest{Type: "poll_local", Offsets: ownedOffsets})
			if err != nil {
				return pollResponse{}, err
			}
			partial = resp.Msgs
		}

		for key, records := range partial {
			msgs[key] = records
		}
	}

	return pollResponse{Type: "poll_ok", Msgs: msgs}, nil
}

func (c Coordinator) CommitOffsets(offsets map[string]int) (commitOffsetsResponse, error) {
	groups := c.groupOffsetsByOwner(offsets)
	for owner, ownedOffsets := range groups {
		if owner == c.node.ID() {
			if err := c.service.CommitOffsets(ownedOffsets); err != nil {
				return commitOffsetsResponse{}, err
			}
			continue
		}

		if _, err := proxy[commitOffsetsResponse](c.node, owner, commitOffsetsLocalRequest{Type: "commit_offsets_local", Offsets: ownedOffsets}); err != nil {
			return commitOffsetsResponse{}, err
		}
	}

	return commitOffsetsResponse{Type: "commit_offsets_ok"}, nil
}

func (c Coordinator) ListCommittedOffsets(keys []string) (listCommittedOffsetsResponse, error) {
	groups := c.groupKeysByOwner(keys)
	offsets := make(map[string]int, len(keys))
	for owner, ownedKeys := range groups {
		var partial map[string]int
		if owner == c.node.ID() {
			partial = c.service.ListCommittedOffsets(ownedKeys)
		} else {
			resp, err := proxy[listCommittedOffsetsResponse](c.node, owner, listCommittedOffsetsLocalRequest{Type: "list_committed_offsets_local", Keys: ownedKeys})
			if err != nil {
				return listCommittedOffsetsResponse{}, err
			}
			partial = resp.Offsets
		}

		for key, offset := range partial {
			offsets[key] = offset
		}
	}

	return listCommittedOffsetsResponse{Type: "list_committed_offsets_ok", Offsets: offsets}, nil
}

func (c Coordinator) groupOffsetsByOwner(offsets map[string]int) map[string]map[string]int {
	grouped := make(map[string]map[string]int)
	for _, key := range sortedMapKeys(offsets) {
		owner := c.service.Owner(key)
		if grouped[owner] == nil {
			grouped[owner] = make(map[string]int)
		}
		grouped[owner][key] = offsets[key]
	}
	return grouped
}

func (c Coordinator) groupKeysByOwner(keys []string) map[string][]string {
	seen := make(map[string]struct{}, len(keys))
	unique := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, key)
	}

	grouped := make(map[string][]string)
	for _, key := range unique {
		owner := c.service.Owner(key)
		grouped[owner] = append(grouped[owner], key)
	}
	return grouped
}

func proxy[T any](node *maelstrom.Node, owner string, body any) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	respMsg, err := node.SyncRPC(ctx, owner, body)
	if err != nil {
		var zero T
		return zero, err
	}

	return maelstromx.Decode[T](respMsg)
}
