package usecase_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/app/usecase"
	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func cartRepositoryWith(cart *domain.Cart) *cartRepositoryStub {
	return &cartRepositoryStub{
		get: func(context.Context, domain.UserID) (*domain.Cart, error) {
			return cart, nil
		},
	}
}

func TestGetCart_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		items []testItemSpec
	}{
		{"empty", nil},
		{"one item", []testItemSpec{{4, 1}}},
		{"several items", []testItemSpec{{4, 1}, {5, 2}, {6, 3}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID := newTestUserID(t, 1)

			products := map[domain.SKU]*domain.Product{}
			for _, spec := range tt.items {
				products[newTestSKU(t, spec.sku)] = newTestProduct(t, spec.sku, uint32(spec.sku)*100)
			}

			var productCalls atomic.Int32
			productProvider := &productProviderStub{
				get: func(_ context.Context, sku domain.SKU) (*domain.Product, error) {
					productCalls.Add(1)
					return products[sku], nil
				},
			}

			var gotUserID domain.UserID
			cartRepository := &cartRepositoryStub{
				get: func(_ context.Context, userID domain.UserID) (*domain.Cart, error) {
					gotUserID = userID
					return newTestCart(t, userID.Value(), tt.items), nil
				},
			}

			got, err := usecase.NewGetCart(cartRepository, productProvider, time.Second).
				Execute(t.Context(), usecase.GetCartRequest{UserID: userID})

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, userID, gotUserID)
			require.Equal(t, userID, got.Cart.UserID())
			require.Equal(t, int32(len(tt.items)), productCalls.Load())

			gotItems := got.Cart.Items()
			require.Len(t, gotItems, len(tt.items))
			for i, spec := range tt.items {
				require.Same(t, products[newTestSKU(t, spec.sku)], gotItems[i].Product())
				require.Equal(t, newTestCount(t, spec.count), gotItems[i].Count())
			}
		})
	}
}

func TestGetCart_Execute_CartRepositoryError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repoErr error
	}{
		{"not found", &port.ErrUserNotFound{Value: newTestUserID(t, 1)}},
		{"storage", port.ErrStorage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cartRepository := &cartRepositoryStub{
				get: func(context.Context, domain.UserID) (*domain.Cart, error) {
					return nil, tt.repoErr
				},
			}

			got, err := usecase.NewGetCart(cartRepository, &productProviderStub{}, time.Second).
				Execute(t.Context(), usecase.GetCartRequest{UserID: newTestUserID(t, 1)})

			require.Nil(t, got)
			require.ErrorIs(t, err, usecase.ErrGetCart)
			require.ErrorIs(t, err, tt.repoErr)
		})
	}
}

func TestGetCart_Execute_ProductProviderError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		providerErr error
	}{
		{"not found", &port.ErrProductNotFound{Value: newTestSKU(t, 4)}},
		{"provider", port.ErrProductProvider},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			productProvider := &productProviderStub{
				get: func(context.Context, domain.SKU) (*domain.Product, error) {
					return nil, tt.providerErr
				},
			}
			cartRepository := cartRepositoryWith(newTestCart(t, 1, []testItemSpec{{4, 1}}))

			got, err := usecase.NewGetCart(cartRepository, productProvider, time.Second).
				Execute(t.Context(), usecase.GetCartRequest{UserID: newTestUserID(t, 1)})

			require.Nil(t, got)
			require.ErrorIs(t, err, usecase.ErrGetCart)
			require.ErrorIs(t, err, tt.providerErr)
		})
	}
}
