package domain_test

import (
	"strings"
	"testing"

	"github.com/2ips/marketplace/cart/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestNewName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"Latin", "Milk", "Milk", nil},
		{"Cyrillic", "Молоко", "Молоко", nil},
		{"Chinese", "牛奶", "牛奶", nil},

		{"max length Latin", strings.Repeat("A", domain.MaxNameLength), strings.Repeat("A", domain.MaxNameLength), nil},
		{"max length Cyrillic", strings.Repeat("Б", domain.MaxNameLength), strings.Repeat("Б", domain.MaxNameLength), nil},
		{"max length Chinese", strings.Repeat("箱", domain.MaxNameLength), strings.Repeat("箱", domain.MaxNameLength), nil},

		{"too long Latin", strings.Repeat("A", domain.MaxNameLength+1), "", &domain.ErrNameTooLong{Length: domain.MaxNameLength + 1}},
		{"too long Cyrillic", strings.Repeat("Б", domain.MaxNameLength+1), "", &domain.ErrNameTooLong{Length: domain.MaxNameLength + 1}},
		{"too long Chinese", strings.Repeat("箱", domain.MaxNameLength+1), "", &domain.ErrNameTooLong{Length: domain.MaxNameLength + 1}},

		{"trims leading and trailing spaces", "  Milk  ", "Milk", nil},
		{"trims tabs and newlines", "\t\nMilk\n\t", "Milk", nil},

		{"keeps inner spaces", "  Whole  milk  ", "Whole  milk", nil},
		{"keeps tabs and newlines", "\t\nWhole\t\nmilk\t\n", "Whole\t\nmilk", nil},

		{"only spaces", "   ", "", &domain.ErrNameEmpty{}},
		{"only tabs and newlines", "\t\n\r", "", &domain.ErrNameEmpty{}},

		{"max length with surrounding spaces", " " + strings.Repeat("A", domain.MaxNameLength) + " ", strings.Repeat("A", domain.MaxNameLength), nil},
		{"max length with surrounding tabs and newlines", "\t\n" + strings.Repeat("A", domain.MaxNameLength) + "\t\n", strings.Repeat("A", domain.MaxNameLength), nil},

		{"too long with surrounding spaces", " " + strings.Repeat("A", domain.MaxNameLength+1) + " ", "", &domain.ErrNameTooLong{Length: domain.MaxNameLength + 1}},
		{"too long with surrounding tabs and newlines", "\t\n" + strings.Repeat("A", domain.MaxNameLength+1) + "\t\n", "", &domain.ErrNameTooLong{Length: domain.MaxNameLength + 1}},

		{"empty", "", "", &domain.ErrNameEmpty{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewName(tt.input)

			if tt.wantErr != nil {
				require.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got.Value())
		})
	}
}
