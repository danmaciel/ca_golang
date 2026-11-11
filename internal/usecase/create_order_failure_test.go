package usecase

import (
	"errors"
	"testing"

	"github.com/danmaciel/ca_golang/internal/dto"
	"github.com/danmaciel/ca_golang/internal/entity"
	"github.com/danmaciel/ca_golang/internal/event"
	"github.com/danmaciel/ca_golang/pkg/events"
)

// faultingOrderRepository allows deterministic failure injection without a DB.
type faultingOrderRepository struct {
	saveErr   error
	saveCalls int
	stored    *entity.Order
}

func (r *faultingOrderRepository) Save(order *entity.Order) error {
	r.saveCalls++
	if r.saveErr != nil {
		return r.saveErr
	}
	copyOfOrder := *order
	r.stored = &copyOfOrder
	return nil
}

func (*faultingOrderRepository) GetAll() ([]entity.Order, error) {
	return nil, nil
}

// eventDispatchProbe exposes only Dispatch. The embedded interface fulfills
// unused methods of EventDispatcherInterface for this focused use-case test.
type eventDispatchProbe struct {
	events.EventDispatcherInterface
	dispatchErr   error
	dispatchCalls int
	lastEvent     events.EventInterface
}

func (p *eventDispatchProbe) Dispatch(ev events.EventInterface) error {
	p.dispatchCalls++
	p.lastEvent = ev
	return p.dispatchErr
}

func TestCreateOrderSaveFailureDoesNotDispatch(t *testing.T) {
	storageFailure := errors.New("simulated storage failure")
	repo := &faultingOrderRepository{saveErr: storageFailure}
	dispatcher := &eventDispatchProbe{}
	usecase := NewCreateOrderUseCase(repo, event.NewOrderCreated(), dispatcher)

	result, err := usecase.Execute(dto.OrderInputDTO{Price: 10, Tax: 2})
	if !errors.Is(err, storageFailure) {
		t.Fatalf("expected storage failure, got %v", err)
	}
	if result != (dto.OrderOutputDTO{}) {
		t.Fatalf("save failure should not return a created order: %+v", result)
	}
	if repo.saveCalls != 1 || repo.stored != nil {
		t.Fatalf("expected one failed save, got calls=%d stored=%+v", repo.saveCalls, repo.stored)
	}
	if dispatcher.dispatchCalls != 0 {
		t.Fatalf("event must not be dispatched when save failed; got %d dispatches", dispatcher.dispatchCalls)
	}
}

func TestCreateOrderDispatchFailureOccursAfterCommittedSave(t *testing.T) {
	dispatchFailure := errors.New("simulated event dispatch failure")
	repo := &faultingOrderRepository{}
	dispatcher := &eventDispatchProbe{dispatchErr: dispatchFailure}
	usecase := NewCreateOrderUseCase(repo, event.NewOrderCreated(), dispatcher)

	result, err := usecase.Execute(dto.OrderInputDTO{Price: 10, Tax: 2})
	if !errors.Is(err, dispatchFailure) {
		t.Fatalf("expected dispatch failure, got %v", err)
	}
	if repo.saveCalls != 1 || repo.stored == nil {
		t.Fatalf("dispatch failure must not undo the completed save; calls=%d stored=%+v", repo.saveCalls, repo.stored)
	}
	if dispatcher.dispatchCalls != 1 {
		t.Fatalf("expected one dispatch attempt, got %d", dispatcher.dispatchCalls)
	}
	if result.ID == "" || result.ID != repo.stored.ID ||
		result.FinalPrice != 12 || result.FinalPrice != repo.stored.FinalPrice {
		t.Fatalf("persisted order must still be represented in the returned output: %+v", result)
	}
	if dispatcher.lastEvent == nil || dispatcher.lastEvent.GetName() != "OrderCreated" {
		t.Fatal("expected OrderCreated event after successful persistence")
	}
	payload, ok := dispatcher.lastEvent.GetPayload().(dto.OrderOutputDTO)
	if !ok || payload != result {
		t.Fatalf("event payload does not match persisted order: %+v", dispatcher.lastEvent.GetPayload())
	}
}

func TestCreateOrderNoRegisteredHandlerIsNotDeliveryConfirmation(t *testing.T) {
	repo := &faultingOrderRepository{}
	dispatcher := events.NewEventDispatcher() // intentionally no subscribers
	usecase := NewCreateOrderUseCase(repo, event.NewOrderCreated(), dispatcher)

	result, err := usecase.Execute(dto.OrderInputDTO{Price: 10, Tax: 2})
	if err != nil || repo.stored == nil || result.ID != repo.stored.ID {
		t.Fatalf("local dispatcher without subscribers returns success; order=%+v err=%v", result, err)
	}
	// This outcome is intentional: Dispatch returns nil when nobody is
	// subscribed. A successful Execute does not prove external delivery.
}
