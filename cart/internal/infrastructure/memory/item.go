package memory

import (
	"fmt"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

type itemModel struct {
	SKU   uint64
	Name  string
	Price uint32
	Count uint32
}

func itemToModel(item *domain.Item) itemModel {
	product := item.Product()

	return itemModel{
		SKU:   product.SKU().Value(),
		Name:  product.Name().Value(),
		Price: product.Price().Value(),
		Count: item.Count().Value(),
	}
}

func itemFromModel(model itemModel) (*domain.Item, error) {
	sku, err := domain.NewSKU(model.SKU)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	name, err := domain.NewName(model.Name)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	price, err := domain.NewPrice(model.Price)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	count, err := domain.NewCount(model.Count)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	product, err := domain.NewProduct(sku, name, price)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	return domain.NewItem(product, count)
}
