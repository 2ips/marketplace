package grpc

import (
	"math"

	"github.com/2ips/marketplace/product/internal/domain"
	pb "github.com/2ips/marketplace/product/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toModel(product *domain.Product) (*pb.Product, error) {
	if product.SKU().Value() > math.MaxInt64 {
		return nil, status.Errorf(codes.Internal, "SKU must be at most %d, got %d", math.MaxInt64, product.SKU().Value())
	}

	if product.Price().Value() > math.MaxInt32 {
		return nil, status.Errorf(codes.Internal, "Price must be at most %d, got %d", math.MaxInt32, product.Price().Value())
	}

	return &pb.Product{
		Sku:   int64(product.SKU().Value()),
		Name:  product.Name().Value(),
		Price: int32(product.Price().Value()),
	}, nil
}
