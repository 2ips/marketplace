package domain_test

import (
	"math"
	"testing"

	"github.com/2ips/marketplace/product/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewSKU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   uint64
		want    uint64
		wantErr error
	}{
		{"common", 4, 4, nil},
		{"almost zero", 1, 1, nil},
		{"max", math.MaxUint64, math.MaxUint64, nil},
		{"zero", 0, 0, &domain.ErrInvalidSKU{Value: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewSKU(tt.input)

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Value())
		})
	}
}
