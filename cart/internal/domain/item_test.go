package domain_test

import (
	"math"
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	t.Parallel()

	sku := newTestSKU(t, 4)
	count := newTestCount(t, 2)

	got, err := domain.NewItem(sku, count)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, sku, got.SKU())
	require.Equal(t, count, got.Count())
}

func TestItem_AddQuantity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		initial uint32
		add     uint32
		want    uint32
		wantErr error
	}{
		{"common", 2, 3, 5, nil},
		{"up to max", math.MaxUint32 - 1, 1, math.MaxUint32, nil},
		{"overflow", math.MaxUint32, 1, math.MaxUint32, &domain.ErrCountOverflow{Total: math.MaxUint32 + 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			item := newTestItem(t, 4, tt.initial)

			err := item.AddQuantity(newTestCount(t, tt.add))

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, item.Count().Value())
		})
	}
}
