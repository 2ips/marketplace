package domain

type Price uint32

func NewPrice(price uint32) (Price, error) {
	return Price(price), nil
}
