package domain

import (
	"fmt"
	"math"
)

type Item struct {
	sku   SKU
	name  Name
	count Count
	price Price
}

type ErrCountOverflow struct {
	Total uint64
}

func (e *ErrCountOverflow) Error() string {
	return fmt.Sprintf("count must be at most %d, got %d", uint32(math.MaxUint32), e.Total)
}

func NewItem(sku SKU, name Name, count Count, price Price) (*Item, error) {
	return &Item{
		sku:   sku,
		name:  name,
		count: count,
		price: price,
	}, nil
}

func (i *Item) SKU() SKU {
	return i.sku
}

func (i *Item) Name() Name {
	return i.name
}

func (i *Item) Count() Count {
	return i.count
}

func (i *Item) Price() Price {
	return i.price
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
