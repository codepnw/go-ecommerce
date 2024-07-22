package orderUsecases

import (
	"math"

	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/orders"
	"github.com/codepnw/go-ecommerce/internal/orders/orderRepositories"
	"github.com/codepnw/go-ecommerce/internal/products/productRepositories"
)

type IOrderUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindAllOrders(req *orders.OrderFilter) *entities.PaginateRes
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

func (u *orderUsecase) FindAllOrders(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.orderRepo.FindAllOrders(req)
	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}
