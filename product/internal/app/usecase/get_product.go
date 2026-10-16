package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
)

var ErrGetProduct = errors.New("get product error")

type GetProduct struct {
	repository port.ProductRepository
}

type GetProductRequest struct {
	SKU domain.SKU
}

type GetProductResponse struct {
	Product *domain.Product
}

func NewGetProduct(repository port.ProductRepository) *GetProduct {
	return &GetProduct{repository: repository}
}

func (uc *GetProduct) Execute(ctx context.Context, req GetProductRequest) (*GetProductResponse, error) {
	product, err := uc.repository.Get(ctx, req.SKU)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetProduct, err)
	}

	return &GetProductResponse{product}, nil
}
