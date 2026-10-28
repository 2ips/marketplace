package memory

import (
	"fmt"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

type itemModel struct {
	SKU   uint64
	Count uint32
}

func itemToModel(item *domain.Item) itemModel {
	return itemModel{
		SKU:   item.SKU().Value(),
		Count: item.Count().Value(),
	}
}

func itemFromModel(model itemModel) (*domain.Item, error) {
	sku, err := domain.NewSKU(model.SKU)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	count, err := domain.NewCount(model.Count)
	if err != nil {
		return nil, fmt.Errorf("%w: from model error: %w", port.ErrStorage, err)
	}

	return domain.NewItem(sku, count)
}
