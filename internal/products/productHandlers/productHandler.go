package productHandlers

import (
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/internal/products"
	"github.com/codepnw/go-ecommerce/internal/products/productUsecases"
	"github.com/gofiber/fiber/v2"
)

type productHandlerErrCode string

const (
	findOneProductErrCode  productHandlerErrCode = "products-001"
	findAllProductsErrCode productHandlerErrCode = "products-002"
)

type IProductHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindAllProducts(c *fiber.Ctx) error
}

type productHnadler struct {
	cfg            config.Config
	productUsecase productUsecases.IProductUsecase
	filesUsecase   filesUsecases.IFilesUsecase
}

func ProductHandler(cfg config.Config, productUsecase productUsecases.IProductUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductHandler {
	return &productHnadler{
		cfg:            cfg,
		productUsecase: productUsecase,
		filesUsecase:   filesUsecase,
	}
}

func (h *productHnadler) FindOneProduct(c *fiber.Ctx) error {
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

func (h *productHnadler) FindAllProducts(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq: &entities.SortReq{},
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
