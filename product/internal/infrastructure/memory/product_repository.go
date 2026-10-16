package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
)

type ProductRepository struct {
	products map[uint64]*domain.Product
}

func NewProductRepository(reader io.Reader) (*ProductRepository, error) {
	var models map[string]*productModel
	if err := json.NewDecoder(reader).Decode(&models); err != nil {
		return nil, fmt.Errorf("%w: load error: %v", port.ErrStorage, err)
	}

	products := make(map[uint64]*domain.Product, len(models))
	for sku, model := range models {
		product, err := fromModel(sku, model)
		if err != nil {
			return nil, fmt.Errorf("%w: parse error: %w", port.ErrStorage, err)
		}

		products[product.SKU().Value()] = product
	}

	return &ProductRepository{products: products}, nil
}

func (r *ProductRepository) Count(ctx context.Context) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, fmt.Errorf("%w: count error: %w", port.ErrStorage, err)
	}

	return uint64(len(r.products)), nil
}

func (r *ProductRepository) Get(ctx context.Context, sku domain.SKU) (*domain.Product, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("%w: get error: %w", port.ErrStorage, err)
	}

	product, ok := r.products[sku.Value()]
	if !ok {
		return nil, fmt.Errorf("%w: get error: %w", port.ErrStorage, &port.ErrProductNotFound{Value: sku})
	}

	return product, nil
}

func (r *ProductRepository) List(ctx context.Context, page port.PageRequest) iter.Seq2[*domain.Product, error] {
	return func(yield func(*domain.Product, error) bool) {
		offset, end, err := page.Bounds(uint64(len(r.products)))
		if err != nil {
			yield(nil, err)
			return
		}

		var i uint64
		for _, product := range r.products {
			if i >= end {
				return
			}

			if i >= offset {
				if err := ctx.Err(); err != nil {
					yield(nil, err)
					return
				}

				if !yield(product, nil) {
					return
				}
			}

			i++
		}
	}
}
