package orderUsecases

import (
	"github.com/codepnw/go-ecommerce/internal/orders"
	"github.com/codepnw/go-ecommerce/internal/orders/orderRepositories"
	"github.com/codepnw/go-ecommerce/internal/products/productRepositories"
)

type IOrderUsecase interface{
	FindOneOrder(orderId string) (*orders.Order, error)
}

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

func (u *orderUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.orderRepo.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}
	return order, nil
}