package middleware

import (
	"database/sql"
	"fmt"
)

type IMiddlewareRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*Role, error)
}

type middlewareRepository struct {
	db *sql.DB
}

func MiddlewareRepository(db *sql.DB) IMiddlewareRepository {
	return &middlewareRepository{db: db}
}

func (r *middlewareRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
		FROM "oauth"
		WHERE "user_id" = $1
		AND "access_token" = $2;
	`
	var check bool
	err := r.db.QueryRow(query, userId, accessToken).Scan(&check)
	
	return err == nil
}

func (r *middlewareRepository) FindRole() ([]*Role, error) {
	query := `
		SELECT
			"id",
			"title"
		FROM "roles"
		ORDER BY "id" DESC;
	`
	roles := make([]*Role, 0)
	err := r.db.QueryRow(query).Scan(&roles)
	if err != nil {
		return nil, fmt.Errorf("role is empty")
	}

	return roles, nil
}
