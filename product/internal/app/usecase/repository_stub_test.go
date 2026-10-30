package usecase_test

import (
	"context"
	"iter"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
)

type repositoryStub struct {
	count func(ctx context.Context) (uint64, error)
	get   func(ctx context.Context, sku domain.SKU) (*domain.Product, error)
	list  func(ctx context.Context, page port.PageRequest) iter.Seq2[*domain.Product, error]
}

func (s *repositoryStub) Count(ctx context.Context) (uint64, error) {
	return s.count(ctx)
}

func (s *repositoryStub) Get(ctx context.Context, sku domain.SKU) (*domain.Product, error) {
	return s.get(ctx, sku)
}

func (s *repositoryStub) List(ctx context.Context, page port.PageRequest) iter.Seq2[*domain.Product, error] {
	return s.list(ctx, page)
}
