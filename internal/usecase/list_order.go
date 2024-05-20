package usecase

import (
	"github.com/danmaciel/ca_golang/internal/dto"
	"github.com/danmaciel/ca_golang/internal/entity"
	"github.com/danmaciel/ca_golang/pkg/events"
)

type ListOrderOutputDTO struct {
	OrderList []dto.OrderOutputDTO
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	ListOrder       events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	ListOrder events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		ListOrder:       ListOrder,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrderUseCase) Execute() (ListOrderOutputDTO, error) {

	listOrders, err := c.OrderRepository.GetAll()

	if err != nil {
		return ListOrderOutputDTO{}, err
	}

	var listDtoOrders []dto.OrderOutputDTO

	for _, order := range listOrders {

		item := dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
		listDtoOrders = append(listDtoOrders, item)
	}

	dtoOut := ListOrderOutputDTO{OrderList: listDtoOrders}
	c.ListOrder.SetPayload(dtoOut)
	c.EventDispatcher.Dispatch(c.ListOrder)

	return dtoOut, nil
}
