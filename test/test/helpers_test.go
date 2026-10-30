//go:build api

package test

import (
	"errors"
	"io"
	"math"
	"testing"

	pb "github.com/2ips/marketplace/product/pkg/api"
	"github.com/2ips/marketplace/test/internal"
	"github.com/stretchr/testify/require"
)

type product struct {
	sku   int64
	name  string
	price int32
}

type pageMetadata struct {
	totalCount int64
	totalPages int64
	page       int64
}

func catalog(t *testing.T) []product {
	t.Helper()

	metadata, products, err := listProducts(t, 1, math.MaxInt32)

	require.NoError(t, err)
	require.NotNil(t, metadata)
	require.Len(t, products, int(metadata.totalCount))

	return products
}

func listProducts(t *testing.T, page int64, pageSize int64) (*pageMetadata, []product, error) {
	t.Helper()

	stream, err := internal.NewProduct(t).ListProducts(t.Context(), &pb.ListProductsRequest{Page: page, PageSize: pageSize})
	require.NoError(t, err)

	var (
		metadata *pageMetadata
		products []product
	)
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return metadata, products, nil
		}
		if err != nil {
			return metadata, products, err
		}

		switch payload := msg.GetPayload().(type) {
		case *pb.ListProductsResponse_PageMetadata:
			require.Nil(t, metadata)
			metadata = &pageMetadata{
				totalCount: payload.PageMetadata.GetTotalCount(),
				totalPages: payload.PageMetadata.GetTotalPages(),
				page:       payload.PageMetadata.GetPage(),
			}
		case *pb.ListProductsResponse_Product:
			require.NotNil(t, metadata)
			products = append(products, fromModel(payload.Product))
		}
	}
}

func fromModel(model *pb.Product) product {
	return product{sku: model.GetSku(), name: model.GetName(), price: model.GetPrice()}
}

func unknownSKU(t *testing.T) int64 {
	t.Helper()

	var sku int64
	for _, p := range catalog(t) {
		sku = max(sku, p.sku)
	}

	return sku + 1
}
