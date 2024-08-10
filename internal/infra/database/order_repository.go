package database

import (
	"database/sql"

	"github.com/lucas4ndrade/FullcycleCleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) List(from, size int64) (os []entity.Order, err error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders LIMIT ? OFFSET ?", size, from)
	if err != nil {
		return
	}

	for rows.Next() {
		var o entity.Order
		if err = rows.Scan(&o.ID, &o.Price, &o.Tax, &o.FinalPrice); err != nil {
			return
		}
		os = append(os, o)
	}
	if err = rows.Close(); err != nil {
		return
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
