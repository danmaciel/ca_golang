package database

import (
	"database/sql"

	"github.com/danmaciel/ca_golang/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	_, err := r.Db.Exec(
		"INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)",
		order.ID, order.Price, order.Tax, order.FinalPrice,
	)
	return err
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("SELECT count(*) FROM orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) GetAll() ([]entity.Order, error) {
	rows, err := r.Db.Query("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
