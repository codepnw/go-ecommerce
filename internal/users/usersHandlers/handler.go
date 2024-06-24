package usersHandlers

import (
	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/users"
	"github.com/codepnw/go-ecommerce/internal/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userErrCode string

const (
	signupCustomerErrCode userErrCode = "users-001"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg     config.Config
	usecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.Config, usecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signupCustomerErrCode),
			err.Error(),
		).Res()
	}

	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signupCustomerErrCode),
			"email pattern is invalid.",
		).Res()
	}

	// insert customer ; error case from user patterns
	result, err := h.usecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signupCustomerErrCode),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signupCustomerErrCode),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signupCustomerErrCode),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
