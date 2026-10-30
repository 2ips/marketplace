package domain_test

import (
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewProduct(t *testing.T) {
	t.Parallel()

	sku, err := domain.NewSKU(4)
	require.NoError(t, err)
	name, err := domain.NewName("Milk")
	require.NoError(t, err)
	price, err := domain.NewPrice(100)
	require.NoError(t, err)

	got, err := domain.NewProduct(sku, name, price)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, sku, got.SKU())
	require.Equal(t, name, got.Name())
	require.Equal(t, price, got.Price())
}
