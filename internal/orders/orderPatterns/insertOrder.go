package orderPatterns

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/codepnw/go-ecommerce/internal/orders"
)

type IInsertOrderBuilder interface {
	initTransaction() error
	insertOrder() error
	insertProductsOrders() error
	getOrderId() string
	commit() error
}

type insertOrderBuilder struct {
	db  *sql.DB
	tx  *sql.Tx
	req *orders.Order
}

func (b *insertOrderBuilder) getOrderId() string { return b.req.Id }

func (b *insertOrderBuilder) initTransaction() error {
	tx, err := b.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	b.tx = tx
	return nil
}

func (b *insertOrderBuilder) insertOrder() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO "orders" (
			"user_id",
			"contact",
			"address",
			"transfer_slip",
			"status"
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "id";
	`

	err := b.tx.QueryRowContext(
		ctx,
		query,
		b.req.UserId,
		b.req.Contact,
		b.req.Address,
		b.req.TransferSlip,
		b.req.Status,
	).Scan(&b.req.Id)

	if err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert order failed: %v", err)
	}
	return nil
}

func (b *insertOrderBuilder) insertProductsOrders() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO "products_orders" (
			"order_id",
			"qty",
			"product"
		)
		VALUES
	`

	values := make([]any, 0)
	lastIndex := 0

	for i := range b.req.Products {
		values = append(
			values,
			b.req.Id,
			b.req.Products[i].Qty,
			b.req.Products[i].Product,
		)

		if i != len(b.req.Products)-1 {
			query += fmt.Sprintf(`	($%d, $%d, $%d),`, lastIndex+1, lastIndex+2, lastIndex+3)
		} else {
			query += fmt.Sprintf(`	($%d, $%d, $%d);`, lastIndex+1, lastIndex+2, lastIndex+3)
		}

		lastIndex += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		values...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert products_orders failed: %v", err)
	}
	return nil
}

func (b *insertOrderBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertOrderBuilder(db *sql.DB, req *orders.Order) IInsertOrderBuilder {
	return &insertOrderBuilder{
		db:  db,
		req: req,
	}
}

type insertOrderEngineer struct {
	builder IInsertOrderBuilder
}

func InsertOrderEngineer(builder IInsertOrderBuilder) *insertOrderEngineer {
	return &insertOrderEngineer{builder: builder}
}

func (en *insertOrderEngineer) InsertOrder() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	if err := en.builder.insertOrder(); err != nil {
		return "", err
	}

	if err := en.builder.insertProductsOrders(); err != nil {
		return "", err
	}

	if err := en.builder.commit(); err != nil {
		return "", err
	}

	return en.builder.getOrderId(), nil
}
