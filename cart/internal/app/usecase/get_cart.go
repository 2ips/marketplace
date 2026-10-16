package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
)

var ErrGetCart = errors.New("get cart error")

type GetCart struct {
	cartRepository port.CartRepository
}

type GetCartRequest struct {
	UserID domain.UserID
}

type GetCartResponse struct {
	Cart *domain.Cart
}

func NewGetCart(cartRepository port.CartRepository) *GetCart {
	return &GetCart{cartRepository: cartRepository}
}

func (uc *GetCart) Execute(ctx context.Context, req GetCartRequest) (*GetCartResponse, error) {
	cart, err := uc.cartRepository.Get(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetCart, err)
	}

	return &GetCartResponse{cart}, nil
}
