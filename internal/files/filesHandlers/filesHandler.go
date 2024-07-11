package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/files"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type filesHandlersErrCode string

const (
	uploadErr filesHandlersErrCode = "files-001"
	deleteErr filesHandlersErrCode = "files-002"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg     config.Config
	usecase filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.Config, usecase filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}
	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// Files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}
	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + filename,
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := h.usecase.UploadToStorage(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

func (h *filesHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	if err := h.usecase.DeleteFileOnStorage(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
