package productUsecases

import "github.com/codepnw/go-ecommerce/internal/products/productRepositories"

type IProductUsecase interface {

}

type productUsecase struct {
	repo productRepositories.IProductRepository
}

func ProductUsecase(repo productRepositories.IProductRepository) IProductUsecase {
	return &productUsecase{
		repo: repo,
	}
}