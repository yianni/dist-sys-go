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

	if !service.Merge([]int{11, 12, 10}) {
		t.Fatal("Merge() = false, want true when new messages are present")
	}

	if service.Merge([]int{10, 11, 12}) {
		t.Fatal("Merge() = true, want false when all messages already exist")
	}

	if got, want := service.Messages(), []int{10, 11, 12}; !slices.Equal(got, want) {
		t.Fatalf("Messages() = %v, want %v", got, want)
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
