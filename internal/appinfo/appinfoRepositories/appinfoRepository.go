package appinfoRepositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/codepnw/go-ecommerce/internal/appinfo"
)

type IAppinfoRepository interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(id int) error
}

type appinfoRepository struct {
	db *sql.DB
}

func AppinfoRepository(db *sql.DB) IAppinfoRepository {
	return &appinfoRepository{db: db}
}

func (r *appinfoRepository) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	query := `
		SELECT
			"id",
			"title"
		FROM "categories"		
	`

	filterValue := make([]any, 0)
	if req.Title != "" {
		query += `WHERE (LOWER("title") LIKE $1)`
		filterValue = append(filterValue, "%"+strings.ToLower(req.Title)+"%")
	}

	query += ";"

	categories := make([]*appinfo.Category, 0)

	rows, err := r.db.Query(query, filterValue...)
	if err != nil {
		return nil, fmt.Errorf("query category failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category appinfo.Category
		if err := rows.Scan(&category.Id, &category.Title); err != nil {
			return nil, fmt.Errorf("scan categories failed: %v", err)
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return categories, nil
}

func (r *appinfoRepository) InsertCategory(req []*appinfo.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO "categories" (
			"title"
		) VALUES
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	valStack := make([]any, 0)
	for i, cat := range req {
		valStack = append(valStack, cat.Title)

		if i != len(req)-1 {
			query += fmt.Sprintf(`($%d),`, i+1)
		} else {
			query += fmt.Sprintf(`($%d)`, i+1)
		}
	}

	query += `RETURNING "id";`

	rows, err := tx.QueryContext(ctx, query, valStack...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert categories failed: %v", err)
	}

	var index int
	for rows.Next() {
		if err := rows.Scan(&req[index].Id); err != nil {
			return fmt.Errorf("scan categories id failed: %d", err)
		}
		index++
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *appinfoRepository) DeleteCategory(id int) error {
	query := `DELETE FROM "categories" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, id); err != nil {
		return fmt.Errorf("delete category failed: %v", err)
	}
	return nil
}
