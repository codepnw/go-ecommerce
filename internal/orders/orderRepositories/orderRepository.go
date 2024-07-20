package orderRepositories

import "database/sql"

type IOrderRepository interface {}

type orderRepository struct {
	db *sql.DB
}

func OrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}