package usersUsecases

import (
	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/users"
	"github.com/codepnw/go-ecommerce/internal/users/usersRepositories"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg  config.Config
	repo usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.Config, repo usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:  cfg,
		repo: repo,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// insert user
	result, err := u.repo.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}