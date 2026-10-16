package domain

type Item struct {
	sku   SKU
	name  Name
	count Count
	price Price
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
