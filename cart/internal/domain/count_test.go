package domain_test

import (
	"math"
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   uint32
		want    uint32
		wantErr error
	}{
		{"common", 4, 4, nil},
		{"almost zero", 1, 1, nil},
		{"max", math.MaxUint32, math.MaxUint32, nil},
		{"zero", 0, 0, &domain.ErrInvalidCount{Value: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewCount(tt.input)

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Value())
		})
	}
}
