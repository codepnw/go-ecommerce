package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/appinfo"
	"github.com/codepnw/go-ecommerce/internal/appinfo/appinfoUsecases"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoErrCode string

const (
	generateApiKeyErrCode appinfoErrCode = "appinfo-001"
	findCategoryErrCode   appinfoErrCode = "appinfo-002"
	insertCategoryErrCode appinfoErrCode = "appinfo-003"
	deleteCategoryErrCode appinfoErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	InsertCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
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

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErrCode),
			err.Error(),
		).Res()
	}

	category, err := h.usecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}

func (h *appinfoHandler) InsertCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(insertCategoryErrCode),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(insertCategoryErrCode),
			"categories are empty",
		).Res()
	}

	if err := h.usecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertCategoryErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) DeleteCategory(c *fiber.Ctx) error {
	id := strings.Trim(c.Params("id"), " ")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteCategoryErrCode),
			"convert id string to int failed",
		).Res()
	}

	if idInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteCategoryErrCode),
			"id must more than zero",
		).Res()
	}

	if err := h.usecase.DeleteCategory(idInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteCategoryErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId string `json:"category_id"`
		}{
			CategoryId: id,
		},
	).Res()
}
