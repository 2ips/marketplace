package domain

import (
	"slices"
)

type Cart struct {
	userID UserID
	items  []*Item
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

func (c *Cart) AddItem(item *Item) error {
	for _, existing := range c.items {
		if existing.sku == item.sku {
			return existing.AddQuantity(item.count)
		}
	}

	c.items = append(c.items, item)
	return nil
}
