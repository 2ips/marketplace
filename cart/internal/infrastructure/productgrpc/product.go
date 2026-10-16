package productgrpc

import (
	"fmt"

	"github.com/2ips/marketplace/cart/internal/domain"
	pb "github.com/2ips/marketplace/product/pkg/api"
)

func fromModel(model *pb.Product) (*domain.Product, error) {
	if model.GetSku() <= 0 {
		return nil, fmt.Errorf("SKU must be greater than zero, got %d", model.GetSku())
	}

	if model.GetPrice() < 0 {
		return nil, fmt.Errorf("price must be non-negative, got %d", model.GetPrice())
	}

	sku, err := domain.NewSKU(uint64(model.GetSku()))
	if err != nil {
		return nil, err
	}

	name, err := domain.NewName(model.GetName())
	if err != nil {
		return nil, err
	}

	price, err := domain.NewPrice(uint32(model.GetPrice()))
	if err != nil {
		return nil, err
	}

	return domain.NewProduct(sku, name, price)
}
