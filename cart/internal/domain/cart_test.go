package domain_test

import (
	"math"
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

type testItemSpec struct {
	sku   uint64
	count uint32
}

func newTestItems(t *testing.T, specs []testItemSpec) []*domain.Item {
	t.Helper()

	items := make([]*domain.Item, 0, len(specs))
	for _, spec := range specs {
		items = append(items, newTestItem(t, spec.sku, spec.count))
	}
	return items
}

func TestNewCart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		userID uint64
		items  []testItemSpec
	}{
		{"empty", 1, []testItemSpec{}},
		{"one item", 1, []testItemSpec{{4, 1}}},
		{"several items", 1, []testItemSpec{{4, 1}, {5, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userID := newTestUserID(t, tt.userID)
			items := newTestItems(t, tt.items)

			got, err := domain.NewCart(userID, items)

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, userID, got.UserID())
			require.Equal(t, items, got.Items())
		})
	}
}

func TestNewCart_CopiesItems(t *testing.T) {
	t.Parallel()

	items := newTestItems(t, []testItemSpec{{4, 1}})
	cart, err := domain.NewCart(newTestUserID(t, 1), items)
	require.NoError(t, err)

	items[0] = newTestItem(t, 9, 9)

	require.Equal(t, newTestItems(t, []testItemSpec{{4, 1}}), cart.Items())
}

func TestCart_AddItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		initial []testItemSpec
		add     testItemSpec
		want    []testItemSpec
		wantErr error
	}{
		{"empty cart", nil, testItemSpec{4, 1}, []testItemSpec{{4, 1}}, nil},
		{"new SKU", []testItemSpec{{4, 1}}, testItemSpec{5, 2}, []testItemSpec{{4, 1}, {5, 2}}, nil},
		{"existing SKU", []testItemSpec{{4, 1}, {5, 2}}, testItemSpec{5, 3}, []testItemSpec{{4, 1}, {5, 5}}, nil},
		{"existing SKU up to max", []testItemSpec{{4, math.MaxUint32 - 1}}, testItemSpec{4, 1}, []testItemSpec{{4, math.MaxUint32}}, nil},
		{"existing SKU overflow", []testItemSpec{{4, math.MaxUint32}}, testItemSpec{4, 1}, []testItemSpec{{4, math.MaxUint32}}, &domain.ErrCountOverflow{Total: math.MaxUint32 + 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cart, err := domain.NewCart(newTestUserID(t, 1), newTestItems(t, tt.initial))
			require.NoError(t, err)

			err = cart.AddItem(newTestItem(t, tt.add.sku, tt.add.count))

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, newTestItems(t, tt.want), cart.Items())
		})
	}
}
