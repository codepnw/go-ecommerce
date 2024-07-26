package orderUsecases

import (
	"fmt"
	"math"

	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/orders"
	"github.com/codepnw/go-ecommerce/internal/orders/orderRepositories"
	"github.com/codepnw/go-ecommerce/internal/products/productRepositories"
)

type IOrderUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindAllOrders(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
	UpdateOrder(req *orders.Order) (*orders.Order, error)
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

func (u *orderUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	// Check product is exists
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}

		prod, err := u.productRepo.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}

		// Set Price
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.orderRepo.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.orderRepo.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *orderUsecase) UpdateOrder(req *orders.Order) (*orders.Order, error) {
	if err := u.orderRepo.UpdateOrder(req); err != nil {
		return nil, err
	}

	order, err := u.orderRepo.FindOneOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return order, nil
}