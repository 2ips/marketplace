package domain_test

import (
	"math"
	"testing"

	"github.com/2ips/marketplace/product/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewPrice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input uint32
		want  uint32
	}{
		{"common", 100, 100},
		{"zero", 0, 0},
		{"max", math.MaxUint32, math.MaxUint32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewPrice(tt.input)

			require.NoError(t, err)
			require.Equal(t, tt.want, got.Value())
		})
	}
}
