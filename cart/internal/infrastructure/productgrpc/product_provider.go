package productgrpc

import (
	"context"
	"fmt"
	"math"

	"github.com/2ips/marketplace/cart/internal/app/port"
	"github.com/2ips/marketplace/cart/internal/domain"
	pb "github.com/2ips/marketplace/product/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductProvider struct {
	client pb.ProductServiceClient
}

func NewProductProvider(client pb.ProductServiceClient) *ProductProvider {
	return &ProductProvider{client: client}
}

func (p *ProductProvider) Get(ctx context.Context, sku domain.SKU) (*domain.Product, error) {
	if sku.Value() > math.MaxInt64 {
		return nil, fmt.Errorf(
			"%w: get error: sku must be at most %d, got %d",
			port.ErrProductProvider,
			math.MaxInt64,
			sku.Value(),
		)
	}

	res, err := p.client.GetProduct(ctx, &pb.GetProductRequest{Sku: int64(sku.Value())})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("%w: %w", port.ErrProductProvider, &port.ErrProductNotFound{Value: sku})
		}

		return nil, fmt.Errorf("%w: %v", port.ErrProductProvider, err)
	}

	product, err := fromModel(res.GetProduct())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", port.ErrProductProvider, err)
	}

	return product, nil
}
