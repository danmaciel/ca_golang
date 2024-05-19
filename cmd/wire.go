//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/danmaciel/clean_arch_golang/internal/entity"
	"github.com/danmaciel/clean_arch_golang/internal/event"
	"github.com/danmaciel/clean_arch_golang/internal/infra/database"
	"github.com/danmaciel/clean_arch_golang/internal/infra/web"
	"github.com/danmaciel/clean_arch_golang/internal/usecase"
	"github.com/danmaciel/clean_arch_golang/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

var setListOrderdEvent = wire.NewSet(
	event.NewListOrder,
	wire.Bind(new(events.EventInterface), new(*event.ListOrder)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setListOrderdEvent,
		usecase.NewListOrderUseCase,
	)
	return &usecase.ListOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
