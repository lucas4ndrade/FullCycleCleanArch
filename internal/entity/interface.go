package entity

type OrderRepositoryInterface interface {
	Save(order *Order) error
	List(from, size int64) ([]Order, error)
}
