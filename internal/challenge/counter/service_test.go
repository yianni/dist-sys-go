package counter

import (
	"context"
	"fmt"
	"testing"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type fakeKV struct {
	values   map[string]int
	casFails map[string]int
}

func newFakeKV() *fakeKV {
	return &fakeKV{
		values:   make(map[string]int),
		casFails: make(map[string]int),
	}
}

func (f *fakeKV) ReadInt(_ context.Context, key string) (int, error) {
	value, ok := f.values[key]
	if !ok {
		return 0, maelstrom.NewRPCError(keyDoesNotExist, "missing key")
	}

	return value, nil
}

func (f *fakeKV) CompareAndSwap(_ context.Context, key string, from, to any, createIfNotExists bool) error {
	if remaining := f.casFails[key]; remaining > 0 {
		f.casFails[key] = remaining - 1
		return maelstrom.NewRPCError(preconditionFailed, "stale value")
	}

	fromInt, ok := from.(int)
	if !ok {
		return fmt.Errorf("from type %T is not int", from)
	}

	toInt, ok := to.(int)
	if !ok {
		return fmt.Errorf("to type %T is not int", to)
	}

	current, exists := f.values[key]
	if !exists {
		if !createIfNotExists || fromInt != 0 {
			return maelstrom.NewRPCError(keyDoesNotExist, "missing key")
		}

		f.values[key] = toInt
		return nil
	}

	if current != fromInt {
		return maelstrom.NewRPCError(preconditionFailed, "stale value")
	}

	f.values[key] = toInt
	return nil
}

func TestServiceAddCreatesAndIncrementsLocalCounter(t *testing.T) {
	t.Parallel()

	kv := newFakeKV()
	service := NewService(
		kv,
		func() string { return "n1" },
		func() []string { return []string{"n1", "n2"} },
	)

	ctx := context.Background()

	if err := service.Add(ctx, 2); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if err := service.Add(ctx, 3); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if got, want := kv.values[counterKey("n1")], 5; got != want {
		t.Fatalf("local counter = %d, want %d", got, want)
	}
}

func TestServiceAddRetriesOnCASConflict(t *testing.T) {
	t.Parallel()

	kv := newFakeKV()
	kv.values[counterKey("n1")] = 4
	kv.casFails[counterKey("n1")] = 1

	service := NewService(
		kv,
		func() string { return "n1" },
		func() []string { return []string{"n1"} },
	)

	if err := service.Add(context.Background(), 6); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if got, want := kv.values[counterKey("n1")], 10; got != want {
		t.Fatalf("local counter = %d, want %d", got, want)
	}
}

func TestServiceReadSumsAllKnownNodeCounters(t *testing.T) {
	t.Parallel()

	kv := newFakeKV()
	kv.values[counterKey("n0")] = 5
	kv.values[counterKey("n2")] = 8

	service := NewService(
		kv,
		func() string { return "n1" },
		func() []string { return []string{"n0", "n1", "n2"} },
	)

	got, err := service.Read(context.Background())
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if want := 13; got != want {
		t.Fatalf("Read() = %d, want %d", got, want)
	}
}
