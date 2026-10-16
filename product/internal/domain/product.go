package domain

type Product struct {
	sku   SKU
	name  Name
	price Price
}

func NewProduct(sku SKU, name Name, price Price) (*Product, error) {
	return &Product{
		sku:   sku,
		name:  name,
		price: price,
	}, nil
}

func (p *Product) SKU() SKU {
	return p.sku
}

func (p *Product) Name() Name {
	return p.name
}

func (p *Product) Price() Price {
	return p.price
}
