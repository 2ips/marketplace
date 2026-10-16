package port

import (
	"context"
	"errors"
	"fmt"

	"github.com/2ips/marketplace/cart/internal/domain"
)

var ErrProductProvider = errors.New("product provider error")

type ErrProductNotFound struct {
	Value domain.SKU
}

type ProductProvider interface {
	Get(ctx context.Context, sku domain.SKU) (*domain.Product, error)
}

func (e *ErrProductNotFound) Error() string {
	return fmt.Sprintf("cannot find product with SKU %d", e.Value.Value())
}
