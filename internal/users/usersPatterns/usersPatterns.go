package usersPatterns

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/codepnw/go-ecommerce/internal/users"
)

type IInsertUser interface {
	Customer() (IInsertUser, error) // return IInsertUser -> Customer().Result()
	Admin() (IInsertUser, error)
	Result() (*users.UserPassport, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sql.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sql.DB, req *users.UserRegisterReq, isAdmin bool) IInsertUser {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

func newCustomer(db *sql.DB, req *users.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newAdmin(db *sql.DB, req *users.UserRegisterReq) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (u *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO "users" (
			"email",
			"password",
			"username",
			"role_id"
		)
		VALUES ($1, $2, $3, 1)
		RETURNING "id";
	`

	if err := u.db.QueryRowContext(
		ctx,
		query,
		u.req.Email,
		u.req.Password,
		u.req.Username,
	).Scan(&u.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return u, nil
}

func (u *userReq) Admin() (IInsertUser, error) {
	return nil, nil
}

func (u *userReq) Result() (*users.UserPassport, error) {
	query := `
		SELECT
			json_build_object(
				'user', "t",
				'token', NULL
			)
		FROM (
			SELECT
				"u"."id",
				"u"."email",
				"u"."username",
				"u"."role_id"
			FROM "users" "u"
			WHERE "u"."id" = $1
		) AS "t"
	`

	data := make([]byte, 0)
	if err := u.db.QueryRow(query, u.id).Scan(&data); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}

	return user, nil
}
