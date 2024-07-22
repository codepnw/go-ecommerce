package orderPatterns

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codepnw/go-ecommerce/internal/orders"
)

type IFindOrderBuilder interface {
	initQuery()
	initCountQuery()
	buildWhereSearch()
	buildWhereStatus()
	buildWhereDate()
	buildSort()
	buildPaginate()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(data []any)
	setLastIndex(n int)
	getDB() *sql.DB
	reset()
}

type findOrderBuilder struct {
	db        *sql.DB
	req       *orders.OrderFilter
	query     string
	values    []any
	lastIndex int
}

func FindOrderBuilder(db *sql.DB, req *orders.OrderFilter) IFindOrderBuilder {
	return &findOrderBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

func (b *findOrderBuilder) initQuery() {
	b.query += `
		SELECT
			array_to_json(array_agg("at"))
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
			WHERE 1 = 1
	`
}

func (b *findOrderBuilder) initCountQuery() {
	b.query += `
		SELECT
			COUNT (*) AS "count"
		FROM "orders" "o"
		WHERE 1 = 1
	`
}

func (b *findOrderBuilder) buildWhereSearch() {
	if b.req.Search == "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		query := fmt.Sprintf(`
			AND (
				LOWER("o"."user_id") LIKE $%d OR
				LOWER("o"."address") LIKE $%d OR
				LOWER("o"."contact") LIKE $%d 
			)`,
			b.lastIndex+1,
			b.lastIndex+2,
			b.lastIndex+3,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildWhereStatus() {
	if b.req.Status != "" {
		b.values = append(b.values, strings.ToLower(b.req.Status))

		query := fmt.Sprintf(`	AND "o"."status" = $%d`, b.lastIndex+1)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildWhereDate() {
	if b.req.StartDate != "" && b.req.EndDate != "" {
		b.values = append(b.values, b.req.StartDate, b.req.EndDate)

		query := fmt.Sprintf(
			`	AND "o"."created_at" BETWEEN DATE($%d) AND ($%d)::DATE + 1`,
			b.lastIndex+1,
			b.lastIndex+2,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildSort() {
	b.values = append(b.values, b.req.OrderBy)
	b.query += fmt.Sprintf(`	ORDER BY $%d %s`, b.lastIndex+1, b.req.Sort)
	b.lastIndex = len(b.values)
}

func (b *findOrderBuilder) buildPaginate() {
	b.values = append(
		b.values,
		(b.req.Page-1)*b.req.Limit,
		b.req.Limit,
	)
	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastIndex+1, b.lastIndex+2)
	b.lastIndex = len(b.values)
}

func (b *findOrderBuilder) closeQuery() {
	b.query += `
		) AS "at"
	`
}

func (b *findOrderBuilder) getQuery() string      { return b.query }
func (b *findOrderBuilder) setQuery(query string) { b.query = query }

func (b *findOrderBuilder) getValues() []any     { return b.values }
func (b *findOrderBuilder) setValues(data []any) { b.values = data }

func (b *findOrderBuilder) setLastIndex(n int) { b.lastIndex = n }
func (b *findOrderBuilder) getDB() *sql.DB     { return b.db }

func (b *findOrderBuilder) reset() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastIndex = 0
}

// Engineer
type findOrderEngineer struct {
	builder IFindOrderBuilder
}

func FindOrderEngineer(builder IFindOrderBuilder) *findOrderEngineer {
	return &findOrderEngineer{builder: builder}
}

func (en *findOrderEngineer) FindOrders() []*orders.Order {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	en.builder.initQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereDate()
	en.builder.buildSort()
	en.builder.buildPaginate()
	en.builder.closeQuery()

	raws := make([]byte, 0)

	rows, err := en.builder.getDB().Query(en.builder.getQuery(), en.builder.getValues()...)
	if err != nil {
		log.Printf("get orders failed: %v", err)
		return make([]*orders.Order, 0)
	}
	defer rows.Close()

	for rows.Next() {
		var raw byte
		if err := rows.Scan(&raw); err != nil {
			log.Printf("scan orders failed: %v", err)
			return make([]*orders.Order, 0)
		}
		raws = append(raws, raw)
	}

	ordersData := make([]*orders.Order, 0)
	if err := json.Unmarshal(raws, &ordersData); err != nil {
		log.Printf("unmarshal orders failed: %v", err)
		return make([]*orders.Order, 0)
	}

	en.builder.reset()
	return ordersData
}

func (en *findOrderEngineer) CountOrders() int {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	en.builder.initCountQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereDate()

	var count int

	rows, err := en.builder.getDB().Query(en.builder.getQuery(), en.builder.getValues()...)
	if err != nil {
		log.Printf("count orders failed: %v", err)
		return 0
	}

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Printf("scan count orders failed: %v", err)
			return 0
		}
	}

	en.builder.reset()
	return count
}
