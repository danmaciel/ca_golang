package events

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// countingHandler records completed callbacks without modifying the event.
type countingHandler struct {
	calls atomic.Int64
}

func (h *countingHandler) Handle(_ EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.calls.Add(1)
}

// holdingHandler blocks until released, proving that the registry lock is not
// held while callbacks run.
type holdingHandler struct {
	started chan struct{}
	release chan struct{}
	calls   atomic.Int64
}

func (h *holdingHandler) Handle(_ EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.calls.Add(1)
	select {
	case h.started <- struct{}{}:
	default:
	}
	<-h.release
}

func TestDispatchUsesSnapshotWhileSubscriptionsChange(t *testing.T) {
	dispatcher := NewEventDispatcher()
	held := &holdingHandler{started: make(chan struct{}, 1), release: make(chan struct{})}
	added := &countingHandler{}
	if err := dispatcher.Register("OrderCreated", held); err != nil {
		t.Fatal(err)
	}

	complete := make(chan error, 1)
	go func() {
		complete <- dispatcher.Dispatch(NewPayloadEvent("OrderCreated", "first"))
	}()

	select {
	case <-held.started:
	case <-time.After(2 * time.Second):
		t.Fatal("first handler did not start")
	}

	// The active dispatch still has the original handler. Changes must
	// complete independently rather than waiting for the callback to finish.
	changed := make(chan error, 1)
	go func() {
		if err := dispatcher.Remove("OrderCreated", held); err != nil {
			changed <- err
			return
		}
		changed <- dispatcher.Register("OrderCreated", added)
	}()
	select {
	case err := <-changed:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("subscription mutation blocked behind running callback")
	}
	close(held.release)
	select {
	case err := <-complete:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("dispatch did not finish")
	}

	if held.calls.Load() != 1 || added.calls.Load() != 0 {
		t.Fatalf("active dispatch must use old snapshot: old=%d new=%d",
			held.calls.Load(), added.calls.Load())
	}
	if err := dispatcher.Dispatch(NewPayloadEvent("OrderCreated", "next")); err != nil {
		t.Fatal(err)
	}
	if held.calls.Load() != 1 || added.calls.Load() != 1 {
		t.Fatalf("later dispatch must use new subscription: old=%d new=%d",
			held.calls.Load(), added.calls.Load())
	}
}

func TestRegistryOperationsRemainRaceFreeWithConcurrentDispatch(t *testing.T) {
	dispatcher := NewEventDispatcher()
	handler := &countingHandler{}
	const operations = 100
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for n := 0; n < operations; n++ {
				switch worker {
				case 0:
					_ = dispatcher.Register("OrderCreated", handler)
				case 1:
					_ = dispatcher.Remove("OrderCreated", handler)
				case 2:
					_ = dispatcher.Has("OrderCreated", handler)
				case 3:
					if err := dispatcher.Dispatch(NewPayloadEvent("OrderCreated", n)); err != nil {
						t.Errorf("unexpected dispatch error: %v", err)
						return
					}
				}
			}
		}(i)
	}
	wg.Wait()
	dispatcher.Clear()
	if dispatcher.Has("OrderCreated", handler) {
		t.Fatal("Clear must remove active subscriptions")
	}
}

func TestClearDuringDispatchDoesNotCancelCapturedHandler(t *testing.T) {
	dispatcher := NewEventDispatcher()
	held := &holdingHandler{started: make(chan struct{}, 1), release: make(chan struct{})}
	if err := dispatcher.Register("event", held); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- dispatcher.Dispatch(NewPayloadEvent("event", nil)) }()
	select {
	case <-held.started:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not start")
	}
	dispatcher.Clear()
	close(held.release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("dispatch did not complete")
	}
	if held.calls.Load() != 1 {
		t.Fatalf("Clear must not revoke callbacks already running: %d", held.calls.Load())
	}
	if err := dispatcher.Dispatch(NewPayloadEvent("event", nil)); err != nil {
		t.Fatal(err)
	}
	if held.calls.Load() != 1 {
		t.Fatal("Clear must take effect on subsequent dispatches")
	}
}
