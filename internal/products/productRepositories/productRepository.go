package productRepositories

import (
	"database/sql"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
)

type IProductRepository interface {
}

type productRepository struct {
	db           *sql.DB
	cfg          config.Config
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductRepository(db *sql.DB, cfg config.Config, filesUsecase filesUsecases.IFilesUsecase) IProductRepository {
	return &productRepository{
		db:  db,
		cfg: cfg,
		filesUsecase: filesUsecase,
	}
}
