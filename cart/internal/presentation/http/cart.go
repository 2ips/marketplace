package http

import (
	"fmt"
	"math"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/2ips/marketplace/cart/pkg/api"
)

func toCartItem(item *domain.Item) (api.CartItem, error) {
	product := item.Product()

	if product.SKU().Value() > math.MaxInt64 {
		return api.CartItem{}, fmt.Errorf("SKU must be at most %d, got %d", math.MaxInt64, product.SKU().Value())
	}

	return api.CartItem{
		Sku:   int64(product.SKU().Value()),
		Name:  product.Name().Value(),
		Price: product.Price().Value(),
		Count: item.Count().Value(),
	}, nil
}
