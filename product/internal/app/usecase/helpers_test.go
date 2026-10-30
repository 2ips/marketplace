package usecase_test

import (
	"testing"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
	"github.com/stretchr/testify/require"
)

func newTestSKU(t *testing.T, value uint64) domain.SKU {
	t.Helper()

	sku, err := domain.NewSKU(value)
	require.NoError(t, err)
	return sku
}

func newTestProduct(t *testing.T, sku uint64) *domain.Product {
	t.Helper()

	name, err := domain.NewName("Milk")
	require.NoError(t, err)
	price, err := domain.NewPrice(100)
	require.NoError(t, err)

	product, err := domain.NewProduct(newTestSKU(t, sku), name, price)
	require.NoError(t, err)
	return product
}

func newTestPage(t *testing.T, page uint64, pageSize uint64) port.PageRequest {
	t.Helper()

	pageRequest, err := port.NewPageRequest(page, pageSize)
	require.NoError(t, err)
	return pageRequest
}
