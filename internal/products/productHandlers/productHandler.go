package productHandlers

import (
	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/internal/products/productUsecases"
)

type IProductHandler interface{}

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
