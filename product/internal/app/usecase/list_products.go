package usecase

import (
	"context"
	"errors"
	"fmt"
	"iter"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
)

var ErrListProducts = errors.New("list products error")

type ListProducts struct {
	repository port.ProductRepository
}

type ListProductsRequest struct {
	Page port.PageRequest
}

type ListProductsResponse struct {
	TotalCount uint64
	TotalPages uint64
	Page       uint64
	Products   iter.Seq2[*domain.Product, error]
}

func NewListProducts(repository port.ProductRepository) *ListProducts {
	return &ListProducts{repository: repository}
}

func (uc *ListProducts) Execute(ctx context.Context, req ListProductsRequest) (*ListProductsResponse, error) {
	totalCount, err := uc.repository.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListProducts, err)
	}

	return &ListProductsResponse{
		TotalCount: totalCount,
		TotalPages: totalPages(totalCount, req.Page.PageSize),
		Page:       req.Page.Page,
		Products:   uc.products(ctx, req.Page),
	}, nil
}

func (uc *ListProducts) products(ctx context.Context, page port.PageRequest) iter.Seq2[*domain.Product, error] {
	return func(yield func(*domain.Product, error) bool) {
		for product, err := range uc.repository.List(ctx, page) {
			if err != nil {
				yield(nil, fmt.Errorf("%w: %w", ErrListProducts, err))
				return
			}

			if !yield(product, nil) {
				return
			}
		}
	}
}

func totalPages(totalCount uint64, pageSize uint64) uint64 {
	if pageSize == 0 {
		return 0
	}

	pages := totalCount / pageSize
	if totalCount%pageSize != 0 {
		pages++
	}

	return pages
}
