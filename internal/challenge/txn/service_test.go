package txn

import (
	"slices"
	"testing"
)

func intPtr(v int) *int { return &v }

func TestServiceApplyReadsOwnWritesInOrder(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	result, writes := service.Apply([]operation{
		{Kind: "w", Key: 1, Value: intPtr(7)},
		{Kind: "r", Key: 1},
		{Kind: "w", Key: 2, Value: intPtr(9)},
	})

	if len(writes) != 2 {
		t.Fatalf("writes len = %d, want 2", len(writes))
	}

	want := []operation{
		{Kind: "w", Key: 1, Value: intPtr(7)},
		{Kind: "r", Key: 1, Value: intPtr(7)},
		{Kind: "w", Key: 2, Value: intPtr(9)},
	}

	if !slices.EqualFunc(result, want, equalOperation) {
		t.Fatalf("Apply() = %#v, want %#v", result, want)
	}
}

func TestServiceMergePrefersNewerVersion(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	service.Merge([]writeState{{Key: 1, Value: 3, Version: version{Counter: 1, NodeID: "n0"}}})
	service.Merge([]writeState{{Key: 1, Value: 2, Version: version{Counter: 1, NodeID: "n0"}}})
	service.Merge([]writeState{{Key: 1, Value: 5, Version: version{Counter: 2, NodeID: "n1"}}})

	snapshot := service.SnapshotWrites()
	if len(snapshot) != 1 {
		t.Fatalf("snapshot len = %d, want 1", len(snapshot))
	}
	if got := snapshot[0].Value; got != 5 {
		t.Fatalf("merged value = %d, want 5", got)
	}

	_, writes := service.Apply([]operation{{Kind: "w", Key: 2, Value: intPtr(8)}})
	if got, want := writes[0].Version.Counter, uint64(3); got != want {
		t.Fatalf("next local version = %d, want %d", got, want)
	}
}

func TestVersionAfter(t *testing.T) {
	t.Parallel()

	if !(version{Counter: 2, NodeID: "n0"}).After(version{Counter: 1, NodeID: "n9"}) {
		t.Fatal("expected higher counter to win")
	}
	if !(version{Counter: 1, NodeID: "n2"}).After(version{Counter: 1, NodeID: "n1"}) {
		t.Fatal("expected node id tie break to win")
	}
}

func equalOperation(left, right operation) bool {
	if left.Kind != right.Kind || left.Key != right.Key {
		return false
	}
	if left.Value == nil || right.Value == nil {
		return left.Value == nil && right.Value == nil
	}
	return *left.Value == *right.Value
}
