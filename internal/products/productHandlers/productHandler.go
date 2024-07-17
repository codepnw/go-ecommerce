package productHandlers

import (
	"fmt"
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/appinfo"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/files"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/internal/products"
	"github.com/codepnw/go-ecommerce/internal/products/productUsecases"
	"github.com/gofiber/fiber/v2"
)

type productHandlerErrCode string

const (
	findOneProductErrCode  productHandlerErrCode = "products-001"
	findAllProductsErrCode productHandlerErrCode = "products-002"
	insertProductsErrCode  productHandlerErrCode = "products-003"
	updateProductsErrCode  productHandlerErrCode = "products-004"
	deleteProductsErrCode  productHandlerErrCode = "products-005"
)

type IProductHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindAllProducts(c *fiber.Ctx) error
	InsertProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productHandler struct {
	cfg            config.Config
	productUsecase productUsecases.IProductUsecase
	filesUsecase   filesUsecases.IFilesUsecase
}

func ProductHandler(cfg config.Config, productUsecase productUsecases.IProductUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductHandler {
	return &productHandler{
		cfg:            cfg,
		productUsecase: productUsecase,
		filesUsecase:   filesUsecase,
	}
}

func (h *productHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productHandler) FindAllProducts(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findAllProductsErrCode),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productUsecase.FindAllProducts(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productHandler) InsertProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductsErrCode),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductsErrCode),
			"category_id is invalid",
		).Res()
	}

	product, err := h.productUsecase.InsertProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductsErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productHandler) UpdateProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Images:   make([]*entities.Image, 0),
		Category: &appinfo.Category{},
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductsErrCode),
			err.Error(),
		).Res()
	}

	req.Id = productId

	product, err := h.productUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductsErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErrCode),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/products/%s", p.FileName),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErrCode),
			err.Error(),
		).Res()
	}

	if err := h.productUsecase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductsErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
