package usecase_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/app/usecase"
	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func newTestAddItemRequest(t *testing.T, userID uint64, spec testItemSpec) usecase.AddItemRequest {
	t.Helper()

	return usecase.AddItemRequest{
		UserID: newTestUserID(t, userID),
		SKU:    newTestSKU(t, spec.sku),
		Count:  newTestCount(t, spec.count),
	}
}

func existingProductProvider(t *testing.T) *productProviderStub {
	t.Helper()

	return &productProviderStub{
		get: func(_ context.Context, sku domain.SKU) (*domain.Product, error) {
			return newTestProduct(t, sku.Value(), 100), nil
		},
	}
}

func TestAddItem_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		initial []testItemSpec
		add     testItemSpec
		want    []testItemSpec
	}{
		{"new cart", nil, testItemSpec{4, 1}, []testItemSpec{{4, 1}}},
		{"new SKU", []testItemSpec{{4, 1}}, testItemSpec{5, 2}, []testItemSpec{{4, 1}, {5, 2}}},
		{"existing SKU", []testItemSpec{{4, 1}, {5, 2}}, testItemSpec{5, 3}, []testItemSpec{{4, 1}, {5, 5}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID := newTestUserID(t, 1)

			var gotProductSKU domain.SKU
			productProvider := &productProviderStub{
				get: func(_ context.Context, sku domain.SKU) (*domain.Product, error) {
					gotProductSKU = sku
					return newTestProduct(t, sku.Value(), 100), nil
				},
			}

			var gotUserID domain.UserID
			var saved *domain.Cart
			cartRepository := &cartRepositoryStub{
				get: func(_ context.Context, userID domain.UserID) (*domain.Cart, error) {
					gotUserID = userID
					if tt.initial == nil {
						return nil, &port.ErrUserNotFound{Value: userID}
					}
					return newTestCart(t, userID.Value(), tt.initial), nil
				},
				save: func(_ context.Context, cart *domain.Cart) error {
					saved = cart
					return nil
				},
			}

			got, err := usecase.NewAddItem(cartRepository, productProvider, time.Second).
				Execute(t.Context(), newTestAddItemRequest(t, 1, tt.add))

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, newTestSKU(t, tt.add.sku), gotProductSKU)
			require.Equal(t, userID, gotUserID)
			require.NotNil(t, saved)
			require.Equal(t, userID, saved.UserID())
			require.Equal(t, newTestItems(t, tt.want), saved.Items())
		})
	}
}

func TestAddItem_Execute_ProductProviderError(t *testing.T) {
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

			got, err := usecase.NewAddItem(&cartRepositoryStub{}, productProvider, time.Second).
				Execute(t.Context(), newTestAddItemRequest(t, 1, testItemSpec{4, 1}))

			require.Nil(t, got)
			require.ErrorIs(t, err, usecase.ErrAddItem)
			require.ErrorIs(t, err, tt.providerErr)
		})
	}
}

func TestAddItem_Execute_CartRepositoryGetError(t *testing.T) {
	t.Parallel()

	cartRepository := &cartRepositoryStub{
		get: func(context.Context, domain.UserID) (*domain.Cart, error) {
			return nil, port.ErrStorage
		},
	}

	got, err := usecase.NewAddItem(cartRepository, existingProductProvider(t), time.Second).
		Execute(t.Context(), newTestAddItemRequest(t, 1, testItemSpec{4, 1}))

	require.Nil(t, got)
	require.ErrorIs(t, err, usecase.ErrAddItem)
	require.ErrorIs(t, err, port.ErrStorage)
}

func TestAddItem_Execute_CountOverflow(t *testing.T) {
	t.Parallel()

	cartRepository := &cartRepositoryStub{
		get: func(_ context.Context, userID domain.UserID) (*domain.Cart, error) {
			return newTestCart(t, userID.Value(), []testItemSpec{{4, math.MaxUint32}}), nil
		},
	}

	got, err := usecase.NewAddItem(cartRepository, existingProductProvider(t), time.Second).
		Execute(t.Context(), newTestAddItemRequest(t, 1, testItemSpec{4, 1}))

	require.Nil(t, got)
	require.ErrorIs(t, err, usecase.ErrAddItem)
	require.ErrorAs(t, err, new(*domain.ErrCountOverflow))
}

func TestAddItem_Execute_SaveError(t *testing.T) {
	t.Parallel()

	cartRepository := &cartRepositoryStub{
		get: func(_ context.Context, userID domain.UserID) (*domain.Cart, error) {
			return newTestCart(t, userID.Value(), nil), nil
		},
		save: func(context.Context, *domain.Cart) error {
			return port.ErrStorage
		},
	}

	got, err := usecase.NewAddItem(cartRepository, existingProductProvider(t), time.Second).
		Execute(t.Context(), newTestAddItemRequest(t, 1, testItemSpec{4, 1}))

	require.Nil(t, got)
	require.ErrorIs(t, err, usecase.ErrAddItem)
	require.ErrorIs(t, err, port.ErrStorage)
}
