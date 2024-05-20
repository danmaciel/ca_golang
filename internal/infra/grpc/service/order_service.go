package service

import (
	"context"

	"github.com/danmaciel/ca_golang/internal/dto"
	"github.com/danmaciel/ca_golang/internal/infra/grpc/pb"
	"github.com/danmaciel/ca_golang/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := dto.OrderInputDTO{
		Price: float32(in.Price),
		Tax:   float32(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) GetListOrder(ctx context.Context, in *pb.GetListOrderRequest) (*pb.ListOrderResponse, error) {

	orders, err := s.ListOrderUseCase.Execute()

	if err != nil {
		return nil, err
	}

	var listResponse []*pb.CreateOrderResponse

	for _, order := range orders.OrderList {
		listResponse = append(listResponse, &pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return &pb.ListOrderResponse{Orders: listResponse}, nil
}
