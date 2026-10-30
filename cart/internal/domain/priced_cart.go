package domain

import (
	"fmt"
	"math"
	"slices"
)

type PricedCart struct {
	userID UserID
	items  []*PricedItem
}

type ErrTotalPriceOverflow struct {
	Total uint64
}

func (e *ErrTotalPriceOverflow) Error() string {
	return fmt.Sprintf("total price must be at most %d, got %d", uint32(math.MaxUint32), e.Total)
}

func NewPricedCart(userID UserID, items []*PricedItem) (*PricedCart, error) {
	return &PricedCart{userID: userID, items: slices.Clone(items)}, nil
}

func (c *PricedCart) UserID() UserID {
	return c.userID
}

func (c *PricedCart) Items() []*PricedItem {
	return slices.Clone(c.items)
}

func (c *PricedCart) TotalPrice() (Price, error) {
	var total uint64
	for _, item := range c.items {
		total += uint64(item.product.price.Value()) * uint64(item.count.Value())
		if total > math.MaxUint32 {
			return Price{}, &ErrTotalPriceOverflow{Total: total}
		}
	}

	return NewPrice(uint32(total))
}
