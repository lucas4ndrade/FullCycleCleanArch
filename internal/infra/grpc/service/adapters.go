package service

import (
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/grpc/pb"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/usecase"
)

// adaptOrderListResponse adapts the order list output into a protocol buffer Order List response
func adaptOrderListResponse(outputDTO usecase.ListOrderOutputDTO) *pb.OrderList {
	protoOrders := []*pb.Order{}
	for _, o := range outputDTO {
		protoOrders = append(protoOrders, &pb.Order{
			Id:         o.ID,
			Price:      float32(o.Price),
			Tax:        float32(o.Tax),
			FinalPrice: float32(o.FinalPrice),
		})
	}

	return &pb.OrderList{Orders: protoOrders}
}

// adaptOrderListRequest adapts the order list GRPC request into a order list input DTO
func adaptOrderListRequest(in *pb.ListOrderRequest) (dto usecase.ListOrderInputDTO) {
	dto = usecase.GetDefaultOrderInputDTO()
	if in.From > 0 {
		dto.From = int64(in.From)
	}
	if in.Size > 0 {
		dto.Size = int64(in.Size)
	}

	return
}
