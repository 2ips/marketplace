package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

type CartRepository struct {
	mu    sync.RWMutex
	carts map[uint64]cartModel
}

func NewCartRepository() *CartRepository {
	return &CartRepository{carts: make(map[uint64]cartModel)}
}

func (r *CartRepository) Get(ctx context.Context, userID domain.UserID) (*domain.Cart, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("%w: get error: %w", port.ErrStorage, err)
	}

	r.mu.RLock()
	model, ok := r.carts[userID.Value()]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: get error: %w", port.ErrStorage, &port.ErrUserNotFound{Value: userID})
	}

	return cartFromModel(userID, model)
}

func (r *CartRepository) Save(ctx context.Context, cart *domain.Cart) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: save error: %w", port.ErrStorage, err)
	}

	model := cartToModel(cart)

	r.mu.Lock()
	r.carts[cart.UserID().Value()] = model
	r.mu.Unlock()

	return nil
}
