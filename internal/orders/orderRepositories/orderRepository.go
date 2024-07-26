package orderRepositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/codepnw/go-ecommerce/internal/orders"
	"github.com/codepnw/go-ecommerce/internal/orders/orderPatterns"
)

type IOrderRepository interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindAllOrders(req *orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.Order) error
}

type orderRepository struct {
	db *sql.DB
}

func OrderRepository(db *sql.DB) IOrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindOneOrder(orderId string) (*orders.Order, error) {
	query := `
		SELECT
			to_jsonb("t")
		FROM (	
			SELECT
				"o"."id",
				"o"."user_id",
				"o"."transfer_slip",
				"o"."status,
				(
					SELECT
						array_to_json(array_agg("pt"))
					FROM (
						SELECT
							"spo"."id",
							"spo"."qty",
							"spo"."product"
						FROM "products_orders" "spo"
						WHERE "spo"."order_id" = "o"."id"
					) AS "pt"
				) AS "products",
				"o"."address",
				"o"."contact",
				(
					SELECT
						SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))
					FROM "products_orders" "po"
					WHERE "po"."order_id" = "o"."id"
				) AS "total_paid",
				"o"."created_at",
				"o"."updated_at"
			FROM "orders" "o"
			WHERE "o"."id" = $1
		) AS "t";
	`

	orderData := &orders.Order{
		TransferSlip: &orders.TransferSlip{},
		Products:     make([]*orders.ProductsOrder, 0),
	}

	raw := make([]byte, 0)

	if err := r.db.QueryRow(query, orderId).Scan(&raw); err != nil {
		return nil, fmt.Errorf("get order failed: %v", err)
	}

	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return orderData, nil
}

func (r *orderRepository) FindAllOrders(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := orderPatterns.FindOrderBuilder(r.db, req)
	engineer := orderPatterns.FindOrderEngineer(builder)
	fmt.Printf("data: %v", engineer.FindOrders())
	return engineer.FindOrders(), engineer.CountOrders()
}

func (r *orderRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := orderPatterns.InsertOrderBuilder(r.db, req)
	orderId, err := orderPatterns.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}
	return orderId, nil
}

func (r *orderRepository) UpdateOrder(req *orders.Order) error {
	query := `UPDATE "orders" SET`
	
	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(values, req.Status)
		queryWhereStack = append(
			queryWhereStack, 
			fmt.Sprintf(`"status" = $%d?`, lastIndex),
		)
		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(values, req.TransferSlip)
		queryWhereStack = append(
			queryWhereStack, 
			fmt.Sprintf(`"transfer_slip" = $%d?`, lastIndex),
		)
		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(`WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack) - 1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}

	query += queryClose

	if _, err := r.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("update order failed: %v", err)
	}

	return nil
}