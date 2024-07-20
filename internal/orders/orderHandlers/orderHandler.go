package orderHandlers

import (
	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/orders/orderUsecases"
)

type IOrderHandler interface {}

type orderHandler struct {
	cfg config.Config
	usecase orderUsecases.IOrderUsecase
}

func OrderHandler(cfg config.Config, usecase orderUsecases.IOrderUsecase) IOrderHandler {
	return orderHandler{
		cfg: cfg,
		usecase: usecase,
	}
}