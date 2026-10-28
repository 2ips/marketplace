package domain

type PricedItem struct {
	product *Product
	count   Count
}

func NewPricedItem(product *Product, count Count) (*PricedItem, error) {
	return &PricedItem{
		product: product,
		count:   count,
	}, nil
}

func (i *PricedItem) Product() *Product {
	return i.product
}

func (i *PricedItem) Count() Count {
	return i.count
}
