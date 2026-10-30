package usecase_test

import (
	"context"
	"testing"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/app/usecase"
	"github.com/2ips/marketplace/product/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestGetProduct_Execute(t *testing.T) {
	t.Parallel()

	sku := newTestSKU(t, 4)
	product := newTestProduct(t, 4)

	var gotSKU domain.SKU
	repository := &repositoryStub{
		get: func(_ context.Context, sku domain.SKU) (*domain.Product, error) {
			gotSKU = sku
			return product, nil
		},
	}

	got, err := usecase.NewGetProduct(repository).Execute(t.Context(), usecase.GetProductRequest{SKU: sku})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Same(t, product, got.Product)
	require.Equal(t, sku, gotSKU)
}

func TestGetProduct_Execute_RepositoryError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		repoErr error
	}{
		{"not found", &port.ErrProductNotFound{Value: newTestSKU(t, 4)}},
		{"storage", port.ErrStorage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := &repositoryStub{
				get: func(context.Context, domain.SKU) (*domain.Product, error) {
					return nil, tt.repoErr
				},
			}

			got, err := usecase.NewGetProduct(repository).Execute(t.Context(), usecase.GetProductRequest{SKU: newTestSKU(t, 4)})

			require.Nil(t, got)
			require.ErrorIs(t, err, usecase.ErrGetProduct)
			require.ErrorIs(t, err, tt.repoErr)
		})
	}
}
