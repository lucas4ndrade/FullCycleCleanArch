package usecase

import (
	"fmt"

	"github.com/lucas4ndrade/FullcycleCleanArch/internal/entity"
	"github.com/lucas4ndrade/FullcycleCleanArch/pkg/events"
)

type CreateOrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input CreateOrderInputDTO) (output OrderOutputDTO, err error) {
	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	err = order.IsValid()
	if err != nil {
		return
	}

	err = order.CalculateFinalPrice()
	if err != nil {
		return
	}

	if err = c.OrderRepository.Save(&order); err != nil {
		return
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.OrderCreated.SetPayload(dto)
	dispatchErr := c.EventDispatcher.Dispatch(c.OrderCreated)
	if dispatchErr != nil {
		fmt.Printf("[ERROR] Failed to dispatch order %s event after order creation, %v", order.ID, err)
	}

	return
}
