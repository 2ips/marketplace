package domain

import (
	"fmt"
	"math"
	"slices"
)

type Cart struct {
	userID UserID
	items  []*Item
}

type ErrTotalPriceOverflow struct {
	Total uint64
}

func (e *ErrTotalPriceOverflow) Error() string {
	return fmt.Sprintf("total price must be at most %d, got %d", uint32(math.MaxUint32), e.Total)
}

func NewCart(userID UserID, items []*Item) (*Cart, error) {
	return &Cart{userID: userID, items: slices.Clone(items)}, nil
}

func (c *Cart) UserID() UserID {
	return c.userID
}

func (c *Cart) Items() []*Item {
	return slices.Clone(c.items)
}

func (c *Cart) TotalPrice() (Price, error) {
	var total uint64
	for _, item := range c.items {
		total += uint64(item.price.Value()) * uint64(item.count.Value())
		if total > math.MaxUint32 {
			return Price{}, &ErrTotalPriceOverflow{Total: total}
		}
	}

	return NewPrice(uint32(total))
}

func (c *Cart) AddItem(item *Item) error {
	for _, existing := range c.items {
		if existing.sku == item.sku {
			return existing.AddQuantity(item.count)
		}
	}

	c.items = append(c.items, item)
	return nil
}
