package port

import (
	"context"
	"errors"
	"fmt"
	"iter"

	"github.com/2ips/marketplace/product/internal/domain"
)

var ErrStorage = errors.New("storage error")

type ErrProductNotFound struct {
	Value domain.SKU
}

type ProductRepository interface {
	Count(ctx context.Context) (uint64, error)
	Get(ctx context.Context, sku domain.SKU) (*domain.Product, error)
	List(ctx context.Context, page PageRequest) iter.Seq2[*domain.Product, error]
}

func (e *ErrProductNotFound) Error() string {
	return fmt.Sprintf("cannot find product with SKU %d", e.Value)
}
