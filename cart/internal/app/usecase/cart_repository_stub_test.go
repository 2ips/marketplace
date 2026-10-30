package usecase_test

import (
	"context"

	"github.com/2ips/marketplace/cart/internal/domain"
)

type cartRepositoryStub struct {
	get  func(ctx context.Context, userID domain.UserID) (*domain.Cart, error)
	save func(ctx context.Context, cart *domain.Cart) error
}

func (s *cartRepositoryStub) Get(ctx context.Context, userID domain.UserID) (*domain.Cart, error) {
	return s.get(ctx, userID)
}

func (s *cartRepositoryStub) Save(ctx context.Context, cart *domain.Cart) error {
	return s.save(ctx, cart)
}
