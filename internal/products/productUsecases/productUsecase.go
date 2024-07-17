package productUsecases

import (
	"math"

	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/products"
	"github.com/codepnw/go-ecommerce/internal/products/productRepositories"
)

type IProductUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindAllProducts(req *products.ProductFilter) *entities.PaginateRes
	InsertProduct(req *products.Product) (*products.Product, error)
}

type productUsecase struct {
	repo productRepositories.IProductRepository
}

func ProductUsecase(repo productRepositories.IProductRepository) IProductUsecase {
	return &productUsecase{
		repo: repo,
	}
}

func (u *productUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.repo.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productUsecase) FindAllProducts(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.repo.FindAllProducts(req)

	return &entities.PaginateRes{
		Data: products,
		Page: req.Page,
		Limit: req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productUsecase) InsertProduct(req *products.Product) (*products.Product, error) {
	product, err := u.repo.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}