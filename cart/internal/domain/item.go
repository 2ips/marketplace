package domain

import (
	"fmt"
	"math"
)

type Item struct {
	product Product
	count   Count
}

type ErrCountOverflow struct {
	Total uint64
}

func (e *ErrCountOverflow) Error() string {
	return fmt.Sprintf("count must be at most %d, got %d", uint32(math.MaxUint32), e.Total)
}

func NewItem(product Product, count Count) (*Item, error) {
	return &Item{
		product: product,
		count:   count,
	}, nil
}

func (i *Item) Product() Product {
	return i.product
}

func (i *Item) Count() Count {
	return i.count
}

func (i *Item) AddQuantity(count Count) error {
	total := uint64(i.count.Value()) + uint64(count.Value())
	if total > math.MaxUint32 {
		return &ErrCountOverflow{Total: total}
	}

	quantity, err := NewCount(uint32(total))
	if err != nil {
		return err
	}

	i.count = quantity
	return nil
}
