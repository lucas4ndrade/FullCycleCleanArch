package usecase

import (
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/entity"
)

type ListOrderInputDTO struct {
	From int64 `json:"from"`
	Size int64 `json:"size"`
}

// GetDefaultOrderInputDTO gets a new order input DTO with default fields
func GetDefaultOrderInputDTO() ListOrderInputDTO {
	return ListOrderInputDTO{
		From: 0,
		Size: 20,
	}
}

type ListOrderOutputDTO []OrderOutputDTO

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute(input ListOrderInputDTO) (output ListOrderOutputDTO, err error) {
	orders, err := c.OrderRepository.List(input.From, input.Size)
	if err != nil {
		return
	}

	output = ListOrderOutputDTO{}
	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		output = append(output, dto)
	}

	return
}
