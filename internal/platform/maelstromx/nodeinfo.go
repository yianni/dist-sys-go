package maelstromx

import "sort"

type NodeIDFunc func() string
type NodeIDsFunc func() []string

func SortedNodeIDs(nodeIDs NodeIDsFunc) []string {
	ids := append([]string(nil), nodeIDs()...)
	sort.Strings(ids)
	return ids
}

func Peers(selfID NodeIDFunc, nodeIDs NodeIDsFunc) []string {
	ids := SortedNodeIDs(nodeIDs)
	peers := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == selfID() {
			continue
		}
		peers = append(peers, id)
	}

	return peers
}
