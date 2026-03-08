package main

import (
	"testing"

	"github.com/stuttgart-things/homerun2-core-catcher/internal/catcher"
)

func TestCatcherInterface(t *testing.T) {
	// Verify MockCatcher satisfies the Catcher interface
	var _ catcher.Catcher = &catcher.MockCatcher{}
}

func TestMockCatcherRunAndShutdown(t *testing.T) {
	mock := &catcher.MockCatcher{}

	done := make(chan struct{})
	go func() {
		mock.Run()
		close(done)
	}()

	mock.Shutdown()
	<-done
}
