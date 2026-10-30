//go:build api

package test

import (
	"math"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/2ips/marketplace/cart/pkg/api"
	"github.com/2ips/marketplace/test/internal"
	"github.com/stretchr/testify/require"
)

func TestGetCart(t *testing.T) {
	t.Parallel()

	client := internal.NewCart(t)
	userID := newUserID()

	var (
		want       []api.CartItem
		totalPrice uint32
	)
	for i, p := range catalog(t) {
		item := api.CartItem{Sku: p.sku, Name: p.name, Price: uint32(p.price), Count: uint32(i + 1)}

		res, err := client.AddItemWithResponse(t.Context(), userID, item.Sku, api.NewCartItem{Count: item.Count})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())

		want = append(want, item)
		totalPrice += item.Price * item.Count
	}

	res, err := client.GetCartWithResponse(t.Context(), userID)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.ElementsMatch(t, want, res.JSON200.Items)
	require.Equal(t, totalPrice, res.JSON200.TotalPrice)
}

func TestAddItem_SameProduct(t *testing.T) {
	t.Parallel()

	client := internal.NewCart(t)
	userID := newUserID()
	p := catalog(t)[0]

	for _, count := range []uint32{2, 3} {
		res, err := client.AddItemWithResponse(t.Context(), userID, p.sku, api.NewCartItem{Count: count})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())
	}

	res, err := client.GetCartWithResponse(t.Context(), userID)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	item := api.CartItem{Sku: p.sku, Name: p.name, Price: uint32(p.price), Count: 5}
	require.Equal(t, []api.CartItem{item}, res.JSON200.Items)
	require.Equal(t, item.Price*item.Count, res.JSON200.TotalPrice)
}

func TestAddItem_Error(t *testing.T) {
	t.Parallel()

	sku := catalog(t)[0].sku

	tests := []struct {
		name       string
		userID     int64
		sku        int64
		count      uint32
		wantStatus int
	}{
		{name: "zero user id", userID: 0, sku: sku, count: 1, wantStatus: http.StatusBadRequest},
		{name: "negative user id", userID: -1, sku: sku, count: 1, wantStatus: http.StatusBadRequest},
		{name: "zero sku", userID: newUserID(), sku: 0, count: 1, wantStatus: http.StatusBadRequest},
		{name: "negative sku", userID: newUserID(), sku: -1, count: 1, wantStatus: http.StatusBadRequest},
		{name: "zero count", userID: newUserID(), sku: sku, count: 0, wantStatus: http.StatusBadRequest},
		{name: "unknown sku", userID: newUserID(), sku: unknownSKU(t), count: 1, wantStatus: http.StatusPreconditionFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := internal.NewCart(t)

			res, err := client.AddItemWithResponse(t.Context(), tt.userID, tt.sku, api.NewCartItem{Count: tt.count})

			require.NoError(t, err)
			require.Equal(t, tt.wantStatus, res.StatusCode())

			if tt.userID > 0 {
				cart, err := client.GetCartWithResponse(t.Context(), tt.userID)
				require.NoError(t, err)
				require.Equal(t, http.StatusNotFound, cart.StatusCode())
			}
		})
	}
}

func TestAddItem_CountOverflow(t *testing.T) {
	t.Parallel()

	client := internal.NewCart(t)
	userID := newUserID()
	sku := catalog(t)[0].sku

	res, err := client.AddItemWithResponse(t.Context(), userID, sku, api.NewCartItem{Count: math.MaxUint32})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())

	res, err = client.AddItemWithResponse(t.Context(), userID, sku, api.NewCartItem{Count: 1})

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode())
}

func TestGetCart_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		userID     int64
		wantStatus int
	}{
		{name: "zero user id", userID: 0, wantStatus: http.StatusBadRequest},
		{name: "negative user id", userID: -1, wantStatus: http.StatusBadRequest},
		{name: "unknown user", userID: newUserID(), wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := internal.NewCart(t).GetCartWithResponse(t.Context(), tt.userID)

			require.NoError(t, err)
			require.Equal(t, tt.wantStatus, res.StatusCode())
			require.Nil(t, res.JSON200)
		})
	}
}

func newUserID() int64 {
	return rand.Int64N(math.MaxInt64) + 1
}
