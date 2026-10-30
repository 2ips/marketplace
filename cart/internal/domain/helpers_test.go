package domain_test

import (
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func newTestSKU(t *testing.T, value uint64) domain.SKU {
	t.Helper()

	sku, err := domain.NewSKU(value)
	require.NoError(t, err)
	return sku
}

func newTestCount(t *testing.T, value uint32) domain.Count {
	t.Helper()

	count, err := domain.NewCount(value)
	require.NoError(t, err)
	return count
}

func newTestUserID(t *testing.T, value uint64) domain.UserID {
	t.Helper()

	userID, err := domain.NewUserID(value)
	require.NoError(t, err)
	return userID
}

func newTestPrice(t *testing.T, value uint32) domain.Price {
	t.Helper()

	price, err := domain.NewPrice(value)
	require.NoError(t, err)
	return price
}

func newTestProduct(t *testing.T, sku uint64, price uint32) *domain.Product {
	t.Helper()

	name, err := domain.NewName("Milk")
	require.NoError(t, err)

	product, err := domain.NewProduct(newTestSKU(t, sku), name, newTestPrice(t, price))
	require.NoError(t, err)
	return product
}

func newTestItem(t *testing.T, sku uint64, count uint32) *domain.Item {
	t.Helper()

	item, err := domain.NewItem(newTestSKU(t, sku), newTestCount(t, count))
	require.NoError(t, err)
	return item
}

func newTestPricedItem(t *testing.T, sku uint64, price uint32, count uint32) *domain.PricedItem {
	t.Helper()

	item, err := domain.NewPricedItem(newTestProduct(t, sku, price), newTestCount(t, count))
	require.NoError(t, err)
	return item
}
