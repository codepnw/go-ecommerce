package orderUsecases

import (
	"github.com/codepnw/go-ecommerce/internal/orders/orderRepositories"
	"github.com/codepnw/go-ecommerce/internal/products/productRepositories"
)

type IOrderUsecase interface{}

type orderUsecase struct {
	orderRepo   orderRepositories.IOrderRepository
	productRepo productRepositories.IProductRepository
}

func OrderUsecase(orderRepo orderRepositories.IOrderRepository, productRepo productRepositories.IProductRepository) IOrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}
