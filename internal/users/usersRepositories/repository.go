package usersRepositories

import (
	"database/sql"
	"fmt"

	"github.com/codepnw/go-ecommerce/internal/users"
	"github.com/codepnw/go-ecommerce/internal/users/usersPatterns"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
}

type usersRepository struct {
	db *sql.DB
}

func UsersRepository(db *sql.DB) IUsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error) {
	result := usersPatterns.InsertUser(r.db, req, false)

	var err error
	if isAdmin {
		result, err = result.Admin()
		if err != nil {
			return nil, err
		}
	} else {
		result, err = result.Customer()
		if err != nil {
			return nil, err
		}
	}

	user, err := result.Result()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
		SELECT 
			"id",
			"email",
			"password",
			"username",
			"role_id"
		FROM "users"
		WHERE "email" = $1;
	`

	user := new(users.UserCredentialCheck)

	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Username, &user.RoleId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
