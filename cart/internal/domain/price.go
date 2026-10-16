package domain

type Price struct {
	value uint32
}

func NewPrice(price uint32) (Price, error) {
	return Price{value: price}, nil
}

func (p Price) Value() uint32 {
	return p.value
}
