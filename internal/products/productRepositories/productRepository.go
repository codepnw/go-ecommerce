package productRepositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/internal/products"
	"github.com/codepnw/go-ecommerce/internal/products/productPatterns"
)

type IProductRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindAllProducts(req *products.ProductFilter) ([]*products.Product, int)
}

type productRepository struct {
	db           *sql.DB
	cfg          config.Config
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductRepository(db *sql.DB, cfg config.Config, filesUsecase filesUsecases.IFilesUsecase) IProductRepository {
	return &productRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
		SELECT
			to_jsonb("t")
		FROM (
			SELECT
				"p"."id",
				"p"."title",
				"p"."description",
				"p"."price",
				(
					SELECT
						to_jsonb("ct")
					FROM (
						SELECT
							"c"."id",
							"c"."title"
						FROM "categories" "c"
						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
						WHERE "pc"."product_id" = "p"."id"
					) AS "ct"
				) AS "category",
				"p"."created_at",
				"p"."updated_at",
				(
					SELECT
						COALESCE(array_to_json(array_agg("it")), '[]'::json)
					FROM (
						SELECT
							"i"."id",
							"i"."filename",
							"i"."url"
						FROM "images" "i"
						WHERE "i"."product_id" = "p"."id"
					) AS "it"
				) AS "images"
			FROM "products" "p"
			WHERE "p"."id" = $1
			LIMIT 1
		) AS "t";
	`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	rows := r.db.QueryRow(query, productId)
	err := rows.Scan(&productBytes)
	if err != nil {
		return nil, fmt.Errorf("get product failed: %v", err)
	}

	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}

	return product, nil
}

func (r *productRepository) FindAllProducts(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productPatterns.FindProductBuilder(r.db, req)
	engineer := productPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	return result, count
}