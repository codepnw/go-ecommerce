package usersUsecases

import (
	"fmt"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/users"
	"github.com/codepnw/go-ecommerce/internal/users/usersRepositories"
	"github.com/codepnw/go-ecommerce/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
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

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// find user
	user, err := u.repo.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	// sign token
	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("can not sign access_token")
	}

	refreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("can not sign refresh_token")
	}

	// set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id: user.Id,
			Email: user.Email,
			Username: user.Username,
			RoleId: user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err = u.repo.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}