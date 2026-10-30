//go:build api

package test

import (
	"testing"

	pb "github.com/2ips/marketplace/product/pkg/api"
	"github.com/2ips/marketplace/test/internal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetProduct(t *testing.T) {
	t.Parallel()

	for _, want := range catalog(t) {
		t.Run(want.name, func(t *testing.T) {
			t.Parallel()

			res, err := internal.NewProduct(t).GetProduct(t.Context(), &pb.GetProductRequest{Sku: want.sku})

			require.NoError(t, err)
			require.Equal(t, want, fromModel(res.GetProduct()))
		})
	}
}

func TestGetProduct_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sku      int64
		wantCode codes.Code
	}{
		{name: "zero sku", sku: 0, wantCode: codes.InvalidArgument},
		{name: "negative sku", sku: -1, wantCode: codes.InvalidArgument},
		{name: "unknown sku", sku: unknownSKU(t), wantCode: codes.NotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := internal.NewProduct(t).GetProduct(t.Context(), &pb.GetProductRequest{Sku: tt.sku})

			require.Equal(t, tt.wantCode, status.Code(err))
			require.Nil(t, res)
		})
	}
}

func TestListProducts(t *testing.T) {
	t.Parallel()

	products := catalog(t)
	total := int64(len(products))

	tests := []struct {
		name         string
		page         int64
		pageSize     int64
		wantMetadata pageMetadata
		wantCount    int
	}{
		{
			name:         "page size equal to catalog",
			page:         1,
			pageSize:     total,
			wantMetadata: pageMetadata{totalCount: total, totalPages: 1, page: 1},
			wantCount:    len(products),
		},
		{
			name:         "page size larger than catalog",
			page:         1,
			pageSize:     total + 1,
			wantMetadata: pageMetadata{totalCount: total, totalPages: 1, page: 1},
			wantCount:    len(products),
		},
		{
			name:         "first page of one product",
			page:         1,
			pageSize:     1,
			wantMetadata: pageMetadata{totalCount: total, totalPages: total, page: 1},
			wantCount:    1,
		},
		{
			name:         "last page of one product",
			page:         total,
			pageSize:     1,
			wantMetadata: pageMetadata{totalCount: total, totalPages: total, page: total},
			wantCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			metadata, got, err := listProducts(t, tt.page, tt.pageSize)

			require.NoError(t, err)
			require.NotNil(t, metadata)
			require.Equal(t, tt.wantMetadata, *metadata)
			require.Len(t, got, tt.wantCount)
			require.Subset(t, products, got)
		})
	}
}

func TestListProducts_Error(t *testing.T) {
	t.Parallel()

	total := int64(len(catalog(t)))

	tests := []struct {
		name     string
		page     int64
		pageSize int64
		wantCode codes.Code
	}{
		{name: "zero page", page: 0, pageSize: 1, wantCode: codes.InvalidArgument},
		{name: "negative page", page: -1, pageSize: 1, wantCode: codes.InvalidArgument},
		{name: "zero page size", page: 1, pageSize: 0, wantCode: codes.InvalidArgument},
		{name: "negative page size", page: 1, pageSize: -1, wantCode: codes.InvalidArgument},
		{name: "page out of range", page: total + 1, pageSize: 1, wantCode: codes.OutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, got, err := listProducts(t, tt.page, tt.pageSize)

			require.Equal(t, tt.wantCode, status.Code(err))
			require.Empty(t, got)
		})
	}
}
