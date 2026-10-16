package grpc

import (
	"context"
	"math"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/app/usecase"
	"github.com/2ips/marketplace/product/internal/domain"
	pb "github.com/2ips/marketplace/product/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer

	getProduct   *usecase.GetProduct
	listProducts *usecase.ListProducts
}

func NewProductServer(getProduct *usecase.GetProduct, listProducts *usecase.ListProducts) *ProductServer {
	return &ProductServer{getProduct: getProduct, listProducts: listProducts}
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if req.GetSku() < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "sku must be positive, got %d", req.GetSku())
	}

	sku, err := domain.NewSKU(uint64(req.GetSku()))
	if err != nil {
		return nil, toGRPCError(err)
	}

	res, err := s.getProduct.Execute(ctx, usecase.GetProductRequest{SKU: sku})
	if err != nil {
		return nil, toGRPCError(err)
	}

	model, err := toModel(res.Product)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{Product: model}, nil
}

func (s *ProductServer) ListProducts(
	req *pb.ListProductsRequest,
	stream grpc.ServerStreamingServer[pb.ListProductsResponse],
) error {
	if req.GetPage() < 0 {
		return status.Errorf(codes.InvalidArgument, "page must be positive, got %d", req.GetPage())
	}

	if req.GetPageSize() < 0 {
		return status.Errorf(codes.InvalidArgument, "page size must be positive, got %d", req.GetPageSize())
	}

	page, err := port.NewPageRequest(uint64(req.GetPage()), uint64(req.GetPageSize()))
	if err != nil {
		return toGRPCError(err)
	}

	res, err := s.listProducts.Execute(stream.Context(), usecase.ListProductsRequest{Page: page})
	if err != nil {
		return toGRPCError(err)
	}

	if res.TotalCount > math.MaxInt64 {
		return status.Errorf(
			codes.Internal,
			"Products count must be at most %d, got %d",
			math.MaxInt64,
			res.TotalCount,
		)
	}

	if res.TotalPages > math.MaxInt64 {
		return status.Errorf(
			codes.Internal,
			"Pages count must be at most %d, got %d",
			math.MaxInt64,
			res.TotalCount,
		)
	}

	if res.Page > math.MaxInt64 {
		return status.Errorf(codes.Internal, "Page must be at most %d, got %d", math.MaxInt64, res.TotalCount)
	}

	if err := stream.Send(&pb.ListProductsResponse{
		Payload: &pb.ListProductsResponse_PageMetadata{PageMetadata: &pb.PageMetadata{
			TotalCount: int64(res.TotalCount),
			TotalPages: int64(res.TotalPages),
			Page:       int64(res.Page),
		}},
	}); err != nil {
		return err
	}

	for product, err := range res.Products {
		if err != nil {
			return toGRPCError(err)
		}

		model, err := toModel(product)
		if err != nil {
			return err
		}

		if err := stream.Send(
			&pb.ListProductsResponse{Payload: &pb.ListProductsResponse_Product{Product: model}},
		); err != nil {
			return err
		}
	}

	return nil
}
