package broadcast

import (
	"slices"
	"testing"
)

func TestServiceAddAndMessages(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n1" },
		func() []string { return []string{"n1"} },
	)

	if !service.Add(7) {
		t.Fatal("Add(7) = false, want true for first insert")
	}

	if service.Add(7) {
		t.Fatal("Add(7) = true, want false for duplicate insert")
	}

	service.Add(3)
	service.Add(9)

	if got, want := service.Messages(), []int{3, 7, 9}; !slices.Equal(got, want) {
		t.Fatalf("Messages() = %v, want %v", got, want)
	}
}

func TestServiceMerge(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n1" },
		func() []string { return []string{"n1", "n2", "n3"} },
	)

	service.Add(10)

	if got, want := service.Merge([]int{11, 12, 10}), []int{11, 12}; !slices.Equal(got, want) {
		t.Fatalf("Merge() = %v, want %v", got, want)
	}

	if got := service.Merge([]int{10, 11, 12}); len(got) != 0 {
		t.Fatalf("Merge() = %v, want no new messages", got)
	}

	if got, want := service.Messages(), []int{10, 11, 12}; !slices.Equal(got, want) {
		t.Fatalf("Messages() = %v, want %v", got, want)
	}
}

func TestServiceMergeIsOrderInsensitiveAndIdempotent(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n1" },
		func() []string { return []string{"n1", "n2", "n3"} },
	)

	service.Add(1)

	first := service.Merge([]int{4, 2, 4, 3, 2})
	second := service.Merge([]int{3, 2, 4, 2, 3})

	if got, want := first, []int{2, 3, 4}; !slices.Equal(got, want) {
		t.Fatalf("first Merge() = %v, want %v", got, want)
	}

	if len(second) != 0 {
		t.Fatalf("second Merge() = %v, want no new messages", second)
	}

	if got, want := service.Messages(), []int{1, 2, 3, 4}; !slices.Equal(got, want) {
		t.Fatalf("Messages() = %v, want %v", got, want)
	}
}

func TestServiceDrainDirtyReturnsOnlyNewMessages(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n1" },
		func() []string { return []string{"n1", "n2"} },
	)

	service.Add(5)
	service.Merge([]int{7, 5, 9})

	if got, want := service.DrainDirty(), []int{5, 7, 9}; !slices.Equal(got, want) {
		t.Fatalf("DrainDirty() = %v, want %v", got, want)
	}

	if got := service.DrainDirty(); len(got) != 0 {
		t.Fatalf("DrainDirty() = %v, want empty after drain", got)
	}
}

func TestServiceConfigureTopologyUsesAssignedNeighbors(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n2" },
		func() []string { return []string{"n1", "n2", "n3", "n4"} },
	)

	service.ConfigureTopology(map[string][]string{
		"n2": {"n4", "n1", "n4"},
	})

	if got, want := service.Peers(), []string{"n1", "n4"}; !slices.Equal(got, want) {
		t.Fatalf("Peers() = %v, want %v", got, want)
	}
}

func TestServicePeersFallsBackToClusterNodes(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n2" },
		func() []string { return []string{"n1", "n2", "n3"} },
	)

	if got, want := service.Peers(), []string{"n1", "n3"}; !slices.Equal(got, want) {
		t.Fatalf("Peers() = %v, want %v", got, want)
	}
}
