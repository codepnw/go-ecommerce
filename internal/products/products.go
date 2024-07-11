package products

import (
	"github.com/codepnw/go-ecommerce/internal/appinfo"
	"github.com/codepnw/go-ecommerce/internal/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Price       float64           `json:"price"`
	Images      []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // title & description
	*entities.PaginationReq
	*entities.SortReq
}
