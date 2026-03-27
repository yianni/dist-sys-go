package kafka

import (
	"slices"
	"testing"
)

func TestServiceSendAssignsMonotonicOffsets(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	offset0 := service.Send("k1", 7)
	offset1 := service.Send("k1", 9)

	if offset0 != 0 || offset1 != 1 {
		t.Fatalf("offsets = %d, %d, want 0, 1", offset0, offset1)
	}
}

func TestServicePollReturnsRecordsFromRequestedOffset(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	service.Send("k1", 11)
	service.Send("k1", 22)
	service.Send("k1", 33)

	msgs := service.Poll(map[string]int{"k1": 1})
	want := []logRecord{{1, 22}, {2, 33}}
	if got := msgs["k1"]; !slices.Equal(got, want) {
		t.Fatalf("Poll() = %v, want %v", got, want)
	}
}

func TestServiceCommitOffsetsOnlyMovesForward(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	service.Send("k1", 11)
	service.Send("k1", 22)
	service.Send("k1", 33)

	if err := service.CommitOffsets(map[string]int{"k1": 1}); err != nil {
		t.Fatalf("CommitOffsets() error = %v", err)
	}
	if err := service.CommitOffsets(map[string]int{"k1": 0}); err != nil {
		t.Fatalf("CommitOffsets() error = %v", err)
	}

	if got := service.ListCommittedOffsets([]string{"k1"})["k1"]; got != 1 {
		t.Fatalf("committed offset = %d, want 1", got)
	}

	if err := service.CommitOffsets(map[string]int{"k1": 2}); err != nil {
		t.Fatalf("CommitOffsets() error = %v", err)
	}
	if got := service.ListCommittedOffsets([]string{"k1"})["k1"]; got != 2 {
		t.Fatalf("committed offset = %d, want 2", got)
	}
}

func TestServiceCommitOffsetsRejectsUnknownStream(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	if err := service.CommitOffsets(map[string]int{"missing": 0}); err == nil {
		t.Fatal("CommitOffsets() error = nil, want error for unknown stream")
	}
}

func TestServiceCommitOffsetsRejectsOffsetsPastStreamEnd(t *testing.T) {
	t.Parallel()

	service := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n0", "n1"} },
	)

	service.Send("k1", 11)

	if err := service.CommitOffsets(map[string]int{"k1": 1}); err == nil {
		t.Fatal("CommitOffsets() error = nil, want error for out-of-range offset")
	}
}

func TestServiceOwnerIsStableAcrossNodeOrder(t *testing.T) {
	t.Parallel()

	serviceA := NewService(
		func() string { return "n0" },
		func() []string { return []string{"n2", "n0", "n1"} },
	)
	serviceB := NewService(
		func() string { return "n1" },
		func() []string { return []string{"n1", "n0", "n2"} },
	)

	if got, want := serviceA.Owner("orders"), serviceB.Owner("orders"); got != want {
		t.Fatalf("owner mismatch: %q != %q", got, want)
	}
}
