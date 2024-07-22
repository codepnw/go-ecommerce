package orderHandlers

import (
	"strings"
	"time"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/orders"
	"github.com/codepnw/go-ecommerce/internal/orders/orderUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlersErrCode string

const (
	findOneOrderErrCode ordersHandlersErrCode = "orders-001"
	findAllOrderErrCode ordersHandlersErrCode = "orders-002"
)

type IOrderHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindAllOrders(c *fiber.Ctx) error
}

type orderHandler struct {
	cfg     config.Config
	usecase orderUsecases.IOrderUsecase
}

func OrderHandler(cfg config.Config, usecase orderUsecases.IOrderUsecase) IOrderHandler {
	return &orderHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *orderHandler) FindOneOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.usecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErrCode),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *orderHandler) FindAllOrders(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findAllOrderErrCode),
			err.Error(),
		).Res()
	}

	// Paginate
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Page = 5
	}

	// Sort
	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}

	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["id"]
	}

	// Date Format YYYY-MM-DD
	dateFormat := "2006-01-02"

	if req.StartDate != "" {
		start, err := time.Parse(dateFormat, req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findAllOrderErrCode),
				"start date is invalid.",
			).Res()
		}
		req.StartDate = start.Format(dateFormat)
	}

	if req.EndDate != "" {
		end, err := time.Parse(dateFormat, req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findAllOrderErrCode),
				"end date is invalid.",
			).Res()
		}
		req.EndDate = end.Format(dateFormat)
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		h.usecase.FindAllOrders(req),
	).Res()
}
