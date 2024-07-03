package middleware

import (
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/pkg/auth"
	"github.com/codepnw/go-ecommerce/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErrCode middlewareHandlersErrCode = "middlware-001"
	jwtAuthErrCode     middlewareHandlersErrCode = "middleware-002"
	paramsCheckErrCode middlewareHandlersErrCode = "middleware-003"
	authorizeErrCode   middlewareHandlersErrCode = "middleware-004"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authotize(expectRoleId ...int) fiber.Handler
}

type middlewareHandler struct {
	cfg     config.Config
	usecase IMiddlewareUsecase
}

func MiddlewareHandler(cfg config.Config, usecase IMiddlewareUsecase) IMiddlewareHandler {
	return &middlewareHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *middlewareHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewareHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErrCode),
			"rotuer not found",
		).Res()
	}
}

func (h *middlewareHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewareHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErrCode),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.usecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErrCode),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewareHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErrCode),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewareHandler) Authotize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErrCode),
				"user_id is not int type",
			).Res()
		}

		// Find Role
		roles, err := h.usecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.StatusInternalServerError,
				string(authorizeErrCode),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErrCode),
			"no permission to access",
		).Res()
	}
}
