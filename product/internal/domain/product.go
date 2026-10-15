package domain

type Product struct {
	SKU   SKU
	Name  Name
	Price Price
}

func NewProduct(sku uint64, name string, price uint32) (*Product, error) {
	productSKU, err := NewSKU(sku)
	if err != nil {
		return nil, err
	}

	productName, err := NewName(name)
	if err != nil {
		return nil, err
	}

	productPrice, err := NewPrice(price)
	if err != nil {
		return nil, err
	}

	return &Product{
		SKU:   productSKU,
		Name:  productName,
		Price: productPrice,
	}, nil
}
