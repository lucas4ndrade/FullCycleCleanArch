package service

import (
	"context"

	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/grpc/pb"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(
	createOrderUseCase usecase.CreateOrderUseCase,
	listOrderUseCase usecase.ListOrderUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (response *pb.Order, err error) {
	dto := usecase.CreateOrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return
	}

	response = &pb.Order{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}
	return
}

func (s *OrderService) ListOrder(ctx context.Context, in *pb.ListOrderRequest) (response *pb.OrderList, err error) {
	dto := adaptOrderListRequest(in)

	output, err := s.ListOrderUseCase.Execute(dto)
	if err != nil {
		return
	}

	response = adaptOrderListResponse(output)
	return
}
