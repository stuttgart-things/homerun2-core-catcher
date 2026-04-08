package catcher

import (
	"testing"
)

func TestRedisCatcherImplementsHydrator(t *testing.T) {
	// Verify RedisCatcher satisfies the Hydrator interface at compile time.
	var _ Hydrator = (*RedisCatcher)(nil)
}

func TestFileCatcherDoesNotImplementHydrator(t *testing.T) {
	// FileCatcher should NOT implement Hydrator — hydration is Redis-only.
	var c Catcher = NewFileCatcher("dummy.json", 0)
	if _, ok := c.(Hydrator); ok {
		t.Fatal("FileCatcher should not implement Hydrator")
	}
}
