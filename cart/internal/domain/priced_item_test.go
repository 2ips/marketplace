package domain_test

import (
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewPricedItem(t *testing.T) {
	t.Parallel()

	product := newTestProduct(t, 4, 100)
	count := newTestCount(t, 2)

	got, err := domain.NewPricedItem(product, count)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Same(t, product, got.Product())
	require.Equal(t, count, got.Count())
}
