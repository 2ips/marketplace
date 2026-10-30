package usecase_test

import (
	"context"
	"iter"
	"testing"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/app/usecase"
	"github.com/2ips/marketplace/product/internal/domain"
	"github.com/stretchr/testify/require"
)

func collectProducts(seq iter.Seq2[*domain.Product, error]) ([]*domain.Product, error) {
	var products []*domain.Product
	for product, err := range seq {
		if err != nil {
			return products, err
		}
		products = append(products, product)
	}
	return products, nil
}

func TestListProducts_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		total     uint64
		page      uint64
		pageSize  uint64
		wantPages uint64
	}{
		{"empty", 0, 1, 10, 0},
		{"less than page", 3, 1, 10, 1},
		{"exact", 20, 1, 10, 2},
		{"remainder", 10, 1, 3, 4},
		{"page size one", 5, 1, 1, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository := &repositoryStub{
				count: func(context.Context) (uint64, error) { return tt.total, nil },
			}

			got, err := usecase.NewListProducts(repository).Execute(t.Context(), usecase.ListProductsRequest{
				Page: newTestPage(t, tt.page, tt.pageSize),
			})

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.total, got.TotalCount)
			require.Equal(t, tt.wantPages, got.TotalPages)
			require.Equal(t, tt.page, got.Page)
		})
	}
}

func TestListProducts_Execute_CountError(t *testing.T) {
	t.Parallel()

	repository := &repositoryStub{
		count: func(context.Context) (uint64, error) { return 0, port.ErrStorage },
	}

	got, err := usecase.NewListProducts(repository).Execute(t.Context(), usecase.ListProductsRequest{
		Page: newTestPage(t, 1, 10),
	})

	require.Nil(t, got)
	require.ErrorIs(t, err, usecase.ErrListProducts)
	require.ErrorIs(t, err, port.ErrStorage)
}

func TestListProducts_Execute_Products(t *testing.T) {
	t.Parallel()

	page := newTestPage(t, 2, 3)
	products := []*domain.Product{newTestProduct(t, 4), newTestProduct(t, 5), newTestProduct(t, 6)}

	var gotPage port.PageRequest
	repository := &repositoryStub{
		count: func(context.Context) (uint64, error) { return 6, nil },
		list: func(_ context.Context, page port.PageRequest) iter.Seq2[*domain.Product, error] {
			gotPage = page
			return func(yield func(*domain.Product, error) bool) {
				for _, product := range products {
					if !yield(product, nil) {
						return
					}
				}
			}
		},
	}

	got, err := usecase.NewListProducts(repository).Execute(t.Context(), usecase.ListProductsRequest{Page: page})
	require.NoError(t, err)

	gotProducts, err := collectProducts(got.Products)

	require.NoError(t, err)
	require.Equal(t, products, gotProducts)
	require.Equal(t, page, gotPage)
}
