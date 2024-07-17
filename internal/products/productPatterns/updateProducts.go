package productPatterns

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codepnw/go-ecommerce/internal/entities"
	"github.com/codepnw/go-ecommerce/internal/files"
	"github.com/codepnw/go-ecommerce/internal/files/filesUsecases"
	"github.com/codepnw/go-ecommerce/internal/products"
)

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategory() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateProductBuilder struct {
	db             *sql.DB
	tx             *sql.Tx
	req            *products.Product
	filesUsecase   filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastIndexStack int
	values         []any
}

func UpdateProductBuilder(db *sql.DB, req *products.Product, filesUsecase filesUsecases.IFilesUsecase) IUpdateProductBuilder {
	return &updateProductBuilder{
		db:           db,
		req:          req,
		filesUsecase: filesUsecase,
		queryFields:  make([]string, 0),
		values:       make([]any, 0),
	}
}

func (b *updateProductBuilder) initTransaction() error {
	tx, err := b.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateProductBuilder) initQuery() {
	b.query += `UPDATE "products" SET`
}

func (b *updateProductBuilder) updateTitleQuery() {
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		b.lastIndexStack = len(b.values)

		b.queryFields = append(
			b.queryFields,
			fmt.Sprintf(`	"title" = $%d`, b.lastIndexStack),
		)
	}
}

func (b *updateProductBuilder) updateDescriptionQuery() {
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastIndexStack = len(b.values)

		b.queryFields = append(
			b.queryFields,
			fmt.Sprintf(`	"description" = $%d`, b.lastIndexStack),
		)
	}
}

func (b *updateProductBuilder) updatePriceQuery() {
	if b.req.Price != 0 {
		b.values = append(b.values, b.req.Price)
		b.lastIndexStack = len(b.values)

		b.queryFields = append(
			b.queryFields,
			fmt.Sprintf(`	"price" = $%d`, b.lastIndexStack),
		)
	}
}

func (b *updateProductBuilder) updateCategory() error {
	if b.req.Category == nil {
		return nil
	}

	if b.req.Category.Id == 0 {
		return nil
	}

	query := `
		UPDATE "products_categories" SET
			"category_id" = $1
		WHERE "product_id" = $2;
	`

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Category.Id,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update products_categories failed: %v", err)
	}

	return nil
}

func (b *updateProductBuilder) insertImages() error {
	query := `
		INSERT INTO "images" (
			"filename",
			"url",
			"product_id"
		)
		VALUES
	`

	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Images {
		valueStack = append(
			valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`	($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`	($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}

	return nil
}

func (b *updateProductBuilder) getOldImages() []*entities.Image {
	query := `
		SELECT
			"id",
			"filename",
			"url"
		FROM "images"
		WHERE "product_id" = $1;
	`
	images := make([]*entities.Image, 0)

	rows, err := b.db.Query(query, b.req.Id)
	if err != nil {
		return make([]*entities.Image, 0)
	}
	defer rows.Close()

	for rows.Next() {
		var image entities.Image
		if err := rows.Scan(&image.Id, &image.FileName, &image.Url); err != nil {
			return make([]*entities.Image, 0)
		}
		images = append(images, &image)
	}

	return images
}

func (b *updateProductBuilder) deleteOldImages() error {
	query := `DELETE FROM "images" WHERE "product_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		delFileReq := make([]*files.DeleteFileReq, 0)
		for _, img := range images {
			delFileReq = append(delFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/products/%s", img.FileName),
			})
		}
		// Delete Images
		if err := b.filesUsecase.DeleteFileOnStorage(delFileReq); err != nil {
			b.tx.Rollback()
			return err
		}
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete images failed: %v", err)
	}

	return nil
}

func (b *updateProductBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastIndexStack = len(b.values)

	b.query += fmt.Sprintf(`	WHERE "id" = $%d`, b.lastIndexStack)
}

func (b *updateProductBuilder) updateProduct() error {
	fmt.Println(len(b.values))
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update product failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) getQueryFields() []string { return b.queryFields }
func (b *updateProductBuilder) getValues() []any         { return b.values }
func (b *updateProductBuilder) getQuery() string         { return b.query }
func (b *updateProductBuilder) setQuery(query string)    { b.query = query }
func (b *updateProductBuilder) getImagesLen() int        { return len(b.req.Images) }

func (b *updateProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

type updateProductEngineer struct {
	builder IUpdateProductBuilder
}

func UpdateProductEngineer(builder IUpdateProductBuilder) *updateProductEngineer {
	return &updateProductEngineer{builder: builder}
}

func (en *updateProductEngineer) sumQueryFields() {
	en.builder.updateTitleQuery()
	en.builder.updateDescriptionQuery()
	en.builder.updatePriceQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields)-1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}
}

func (en *updateProductEngineer) UpdateProduct() error {
	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	// Update Product
	if err := en.builder.updateProduct(); err != nil {
		return err
	}

	// Update Category
	if err := en.builder.updateCategory(); err != nil {
		return err
	}

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.deleteOldImages(); err != nil {
			return err
		}
		if err := en.builder.insertImages(); err != nil {
			return err
		}
	}

	// Commit
	if err := en.builder.commit(); err != nil {
		return err
	}
	return nil
}
