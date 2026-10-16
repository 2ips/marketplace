package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

var ErrAddItem = errors.New("add item error")

type AddItem struct {
	cartRepository  port.CartRepository
	productProvider port.ProductProvider
}

type AddItemRequest struct {
	UserID domain.UserID
	SKU    domain.SKU
	Count  domain.Count
}

type AddItemResponse struct{}

func NewAddItem(cartRepository port.CartRepository, productProvider port.ProductProvider) *AddItem {
	return &AddItem{cartRepository: cartRepository, productProvider: productProvider}
}

func (uc *AddItem) Execute(ctx context.Context, req AddItemRequest) (*AddItemResponse, error) {
	product, err := uc.productProvider.Get(ctx, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAddItem, err)
	}

	item, err := domain.NewItem(product, req.Count)
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
