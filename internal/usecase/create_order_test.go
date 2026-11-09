package usecase

import (
	"sync"
	"testing"

	"github.com/danmaciel/ca_golang/internal/dto"
	"github.com/danmaciel/ca_golang/internal/entity"
	"github.com/danmaciel/ca_golang/internal/event"
	"github.com/danmaciel/ca_golang/pkg/events"
)

type recordingOrders struct {
	mu    sync.Mutex
	saved int
}

func (r *recordingOrders) Save(order *entity.Order) error {
	r.mu.Lock()
	r.saved++
	r.mu.Unlock()
	return nil
}

func (*recordingOrders) GetAll() ([]entity.Order, error) { return nil, nil }

type capturedOrderHandler struct {
	payloads chan dto.OrderOutputDTO
}

func (h *capturedOrderHandler) Handle(ev events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.payloads <- ev.GetPayload().(dto.OrderOutputDTO)
}

func TestCreateOrderRejectsInvalidDataWithoutSavingOrDispatch(t *testing.T) {
	repo := &recordingOrders{}
	dispatcher := events.NewEventDispatcher()
	handler := &capturedOrderHandler{payloads: make(chan dto.OrderOutputDTO, 1)}
	if err := dispatcher.Register("OrderCreated", handler); err != nil {
		t.Fatal(err)
	}
	usecase := NewCreateOrderUseCase(repo, event.NewOrderCreated(), dispatcher)
	_, err := usecase.Execute(dto.OrderInputDTO{Price: 0, Tax: 2})
	if err == nil || repo.saved != 0 || len(handler.payloads) != 0 {
		t.Fatalf("invalid order should neither persist nor dispatch: err=%v saved=%d events=%d", err, repo.saved, len(handler.payloads))
	}
}

func TestConcurrentCreatesDispatchIndependentPayloads(t *testing.T) {
	const total = 24
	repo := &recordingOrders{}
	dispatcher := events.NewEventDispatcher()
	handler := &capturedOrderHandler{payloads: make(chan dto.OrderOutputDTO, total)}
	if err := dispatcher.Register("OrderCreated", handler); err != nil {
		t.Fatal(err)
	}
	template := event.NewOrderCreated()
	usecase := NewCreateOrderUseCase(repo, template, dispatcher)

	var wg sync.WaitGroup
	failures := make(chan error, total)
	for i := 1; i <= total; i++ {
		wg.Add(1)
		go func(price int) {
			defer wg.Done()
			_, err := usecase.Execute(dto.OrderInputDTO{Price: float32(price), Tax: 1})
			if err != nil {
				failures <- err
			}
		}(i)
	}
	wg.Wait()
	close(failures)
	for err := range failures {
		t.Errorf("create failed: %v", err)
	}

	if repo.saved != total || len(handler.payloads) != total {
		t.Fatalf("expected %d saves and events, got %d and %d", total, repo.saved, len(handler.payloads))
	}
	seen := make(map[float32]bool)
	for i := 0; i < total; i++ {
		payload := <-handler.payloads
		if payload.FinalPrice != payload.Price+payload.Tax {
			t.Errorf("wrong final price in payload: %+v", payload)
		}
		if seen[payload.Price] {
			t.Errorf("duplicate event payload price %v", payload.Price)
		}
		seen[payload.Price] = true
	}
	if template.GetPayload() != nil {
		t.Fatal("shared injected event template was mutated")
	}
}
