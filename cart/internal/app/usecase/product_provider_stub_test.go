package usecase_test

import (
	"context"

	"github.com/2ips/marketplace/cart/internal/domain"
)

type productProviderStub struct {
	get func(ctx context.Context, sku domain.SKU) (*domain.Product, error)
}

func (s *productProviderStub) Get(ctx context.Context, sku domain.SKU) (*domain.Product, error) {
	return s.get(ctx, sku)
}
