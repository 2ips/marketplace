package domain

import "fmt"

type SKU uint64

type ErrInvalidSKU struct {
	Value uint64
}

func (e *ErrInvalidSKU) Error() string {
	return fmt.Sprintf("SKU must be non zero, got %d", e.Value)
}

func NewSKU(sku uint64) (SKU, error) {
	if sku == 0 {
		return 0, &ErrInvalidSKU{Value: sku}
	}

	return SKU(sku), nil
}
