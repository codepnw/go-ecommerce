package appinfoHandlers

import (
	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoUsecases"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoErrCode string

const (
	generateApiKeyErrCode appinfoErrCode = "appinfo-001"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg     config.Config
	usecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.Config, usecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth(
		auth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}
