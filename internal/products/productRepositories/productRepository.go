package productRepositories

import (
	"context"
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
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
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

func (r *productRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builer := productPatterns.InsertProductBuilder(r.db, req)

	productId, err := productPatterns.InsertProductEngineer(builer).InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productPatterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := productPatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) DeleteProduct(productId string) error {
	query := `DELETE "products" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, productId); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}

	return nil
}