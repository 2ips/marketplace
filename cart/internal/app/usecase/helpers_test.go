package usecase_test

import (
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

type testItemSpec struct {
	sku   uint64
	count uint32
}

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

func newTestProduct(t *testing.T, sku uint64, price uint32) *domain.Product {
	t.Helper()

	name, err := domain.NewName("Milk")
	require.NoError(t, err)
	productPrice, err := domain.NewPrice(price)
	require.NoError(t, err)

	product, err := domain.NewProduct(newTestSKU(t, sku), name, productPrice)
	require.NoError(t, err)
	return product
}

func newTestItems(t *testing.T, specs []testItemSpec) []*domain.Item {
	t.Helper()

	items := make([]*domain.Item, 0, len(specs))
	for _, spec := range specs {
		item, err := domain.NewItem(newTestSKU(t, spec.sku), newTestCount(t, spec.count))
		require.NoError(t, err)
		items = append(items, item)
	}
	return items
}

func newTestCart(t *testing.T, userID uint64, specs []testItemSpec) *domain.Cart {
	t.Helper()

	cart, err := domain.NewCart(newTestUserID(t, userID), newTestItems(t, specs))
	require.NoError(t, err)
	return cart
}
