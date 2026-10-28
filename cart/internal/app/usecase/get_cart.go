package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
	"golang.org/x/sync/errgroup"
)

var ErrGetCart = errors.New("get cart error")

type GetCart struct {
	cartRepository  port.CartRepository
	productProvider port.ProductProvider
	productTimeout  time.Duration
}

type GetCartRequest struct {
	UserID domain.UserID
}

type GetCartResponse struct {
	Cart *domain.PricedCart
}

func NewGetCart(
	cartRepository port.CartRepository,
	productProvider port.ProductProvider,
	productTimeout time.Duration,
) *GetCart {
	return &GetCart{
		cartRepository:  cartRepository,
		productProvider: productProvider,
		productTimeout:  productTimeout,
	}
}

func (uc *GetCart) Execute(ctx context.Context, req GetCartRequest) (*GetCartResponse, error) {
	cart, err := uc.cartRepository.Get(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetCart, err)
	}

	items := cart.Items()

	var products []*domain.Product
	if len(items) > 0 {
		products, err = uc.getProducts(ctx, items)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrGetCart, err)
		}
	}

	pricedItems := make([]*domain.PricedItem, len(items))
	for i, item := range items {
		pricedItem, err := domain.NewPricedItem(products[i], item.Count())
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrGetCart, err)
		}

		pricedItems[i] = pricedItem
	}

	pricedCart, err := domain.NewPricedCart(cart.UserID(), pricedItems)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetCart, err)
	}

	return &GetCartResponse{pricedCart}, nil
}

func (uc *GetCart) getProducts(ctx context.Context, items []*domain.Item) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.productTimeout)
	defer cancel()

	products := make([]*domain.Product, len(items))

	group, ctx := errgroup.WithContext(ctx)
	for i, item := range items {
		group.Go(func() error {
			product, err := uc.productProvider.Get(ctx, item.SKU())
			if err != nil {
				return err
			}

			products[i] = product
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return products, nil
}
