package domain

type Product struct {
	SKU   SKU
	Name  Name
	Price Price
}

func NewProduct(sku SKU, name Name, price Price) (*Product, error) {
	return &Product{
		SKU:   sku,
		Name:  name,
		Price: price,
	}, nil
}
