package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

var ErrAddItem = errors.New("add item error")

type AddItem struct {
	cartRepository  port.CartRepository
	productProvider port.ProductProvider
	productTimeout  time.Duration
}

type AddItemRequest struct {
	UserID domain.UserID
	SKU    domain.SKU
	Count  domain.Count
}

type AddItemResponse struct{}

func NewAddItem(
	cartRepository port.CartRepository,
	productProvider port.ProductProvider,
	productTimeout time.Duration,
) *AddItem {
	return &AddItem{
		cartRepository:  cartRepository,
		productProvider: productProvider,
		productTimeout:  productTimeout,
	}
}

func (uc *AddItem) Execute(ctx context.Context, req AddItemRequest) (*AddItemResponse, error) {
	productCtx, cancel := context.WithTimeout(ctx, uc.productTimeout)
	defer cancel()

	// TODO: add product existence check method to product service
	if _, err := uc.productProvider.Get(productCtx, req.SKU); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
	}

	item, err := domain.NewItem(req.SKU, req.Count)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
	}

	cart, err := uc.cartRepository.Get(ctx, req.UserID)
	if err != nil {
		if _, ok := errors.AsType[*port.ErrUserNotFound](err); !ok {
			return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
		}

		cart, err = domain.NewCart(req.UserID, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
		}
	}

	if err := cart.AddItem(item); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
	}

	if err := uc.cartRepository.Save(ctx, cart); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
	}

	return &AddItemResponse{}, nil
}
