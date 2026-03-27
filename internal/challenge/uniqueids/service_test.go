package uniqueids

import (
	"strings"
	"sync"
	"testing"
)

func TestServiceGenerateSequential(t *testing.T) {
	t.Parallel()

	service := NewService(func() string { return "n1" })

	seen := make(map[string]struct{}, 1000)
	for range 1000 {
		id := service.Generate()
		if !strings.HasPrefix(id, "n1-") {
			t.Fatalf("Generate() prefix = %q, want prefix %q", id, "n1-")
		}

		if _, exists := seen[id]; exists {
			t.Fatalf("Generate() produced duplicate id %q", id)
		}

		seen[id] = struct{}{}
	}
}

func TestServiceGenerateConcurrent(t *testing.T) {
	t.Parallel()

	service := NewService(func() string { return "n2" })

	const goroutines = 8
	const perGoroutine = 250

	results := make(chan string, goroutines*perGoroutine)

	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range perGoroutine {
				results <- service.Generate()
			}
		}()
	}

	wg.Wait()
	close(results)

	seen := make(map[string]struct{}, goroutines*perGoroutine)
	for id := range results {
		if !strings.HasPrefix(id, "n2-") {
			t.Fatalf("Generate() prefix = %q, want prefix %q", id, "n2-")
		}

		if _, exists := seen[id]; exists {
			t.Fatalf("Generate() produced duplicate id %q", id)
		}

		seen[id] = struct{}{}
	}

	if got, want := len(seen), goroutines*perGoroutine; got != want {
		t.Fatalf("generated %d unique ids, want %d", got, want)
	}
}
