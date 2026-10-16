package memory

import (
	"fmt"
	"strconv"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
)

type productModel struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func fromModel(sku string, product *productModel) (*domain.Product, error) {
	skuValue, err := strconv.ParseUint(sku, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %v", port.ErrStorage, err)
	}

	domainSKU, err := domain.NewSKU(skuValue)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	name, err := domain.NewName(product.Name)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	price, err := domain.NewPrice(product.Price)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	return domain.NewProduct(domainSKU, name, price)
}
