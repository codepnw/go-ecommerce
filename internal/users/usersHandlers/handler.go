package usersHandlers

import (
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/users"
	"github.com/codepnw/go-ecommerce/internal/users/usersUsecases"
	"github.com/codepnw/go-ecommerce/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type userErrCode string

const (
	signupCustomerErrCode     userErrCode = "users-001"
	signInErrCode             userErrCode = "users-002"
	refreshErrCode            userErrCode = "users-003"
	signoutErrCode            userErrCode = "users-004"
	signupAdminErrCode        userErrCode = "users-005"
	generateAdminTokenErrCode userErrCode = "users-006"
	getUserProfileErrCode     userErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
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

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
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
	result, err := h.usecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signupAdminErrCode),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signupAdminErrCode),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signupAdminErrCode),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	passport, err := h.usecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshErrCode),
			err.Error(),
		).Res()
	}

	passport, err := h.usecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	if err := h.usecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewAuth(auth.Admin, h.cfg.Jwt(), nil)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusInternalServerError,
			string(generateAdminTokenErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	// Get profile
	result, err := h.usecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErrCode),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErrCode),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
