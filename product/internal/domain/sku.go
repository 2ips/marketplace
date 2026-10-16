package domain

import "fmt"

type SKU struct {
	value uint64
}

type ErrInvalidSKU struct {
	Value uint64
}

func (e *ErrInvalidSKU) Error() string {
	return fmt.Sprintf("SKU must be non zero, got %d", e.Value)
}

func NewSKU(sku uint64) (SKU, error) {
	if sku == 0 {
		return SKU{}, &ErrInvalidSKU{Value: sku}
	}

	return SKU{value: sku}, nil
}

func (s SKU) Value() uint64 {
	return s.value
}
