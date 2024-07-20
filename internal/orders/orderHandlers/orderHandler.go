package orderHandlers

import (
	"strings"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/orders/orderUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlersErrCode string

const (
	findOneOrderErrCode ordersHandlersErrCode = "orders-001"
)

type IOrderHandler interface{
	FindOneOrder(c *fiber.Ctx) error
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
