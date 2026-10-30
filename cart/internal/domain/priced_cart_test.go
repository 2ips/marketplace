package domain_test

import (
	"math"
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

type testPricedItemSpec struct {
	price uint32
	count uint32
}

func newTestPricedItems(t *testing.T, specs []testPricedItemSpec) []*domain.PricedItem {
	t.Helper()

	items := make([]*domain.PricedItem, 0, len(specs))
	for i, spec := range specs {
		items = append(items, newTestPricedItem(t, uint64(i+1), spec.price, spec.count))
	}
	return items
}

func TestNewPricedCart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		userID uint64
		items  []testPricedItemSpec
	}{
		{"empty", 1, []testPricedItemSpec{}},
		{"one item", 1, []testPricedItemSpec{{100, 1}}},
		{"several items", 1, []testPricedItemSpec{{100, 1}, {200, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID := newTestUserID(t, tt.userID)
			items := newTestPricedItems(t, tt.items)

			got, err := domain.NewPricedCart(userID, items)

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, userID, got.UserID())
			require.Equal(t, items, got.Items())
		})
	}
}

func TestNewPricedCart_CopiesItems(t *testing.T) {
	t.Parallel()

	items := newTestPricedItems(t, []testPricedItemSpec{{100, 1}})
	cart, err := domain.NewPricedCart(newTestUserID(t, 1), items)
	require.NoError(t, err)

	items[0] = newTestPricedItem(t, 9, 999, 9)

	require.Equal(t, newTestPricedItems(t, []testPricedItemSpec{{100, 1}}), cart.Items())
}

func TestPricedCart_TotalPrice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		items   []testPricedItemSpec
		want    uint32
		wantErr error
	}{
		{"empty cart", nil, 0, nil},
		{"one item", []testPricedItemSpec{{100, 3}}, 300, nil},
		{"several items", []testPricedItemSpec{{100, 3}, {50, 2}, {1, 1}}, 401, nil},
		{"zero price item", []testPricedItemSpec{{0, 5}, {100, 1}}, 100, nil},
		{"exactly max", []testPricedItemSpec{{math.MaxUint32 - 1, 1}, {1, 1}}, math.MaxUint32, nil},
		{"sum overflow", []testPricedItemSpec{{math.MaxUint32, 1}, {1, 1}}, 0, &domain.ErrTotalPriceOverflow{Total: math.MaxUint32 + 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cart, err := domain.NewPricedCart(newTestUserID(t, 1), newTestPricedItems(t, tt.items))
			require.NoError(t, err)

			got, err := cart.TotalPrice()

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Value())
		})
	}
}
