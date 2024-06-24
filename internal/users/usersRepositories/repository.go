package usersRepositories

import (
	"database/sql"

	"github.com/codepnw/go-ecommerce/internal/users"
	"github.com/codepnw/go-ecommerce/internal/users/usersPatterns"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq, isAdmin bool) (*users.UserPassport, error)
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