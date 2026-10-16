package grpc

import (
	"errors"
	"fmt"

	"github.com/2ips/marketplace/product/internal/app/port"
	"github.com/2ips/marketplace/product/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCError(err error) error {
	var (
		notFound    *port.ErrProductNotFound
		invalidSKU  *domain.ErrInvalidSKU
		invalidPage *port.ErrInvalidPage
		invalidSize *port.ErrInvalidPageSize
		outOfRange  *port.ErrPageOutOfRange
	)

	switch {
	case errors.As(err, &notFound):
		return status.Error(codes.NotFound, notFound.Error())
	case errors.As(err, &invalidSKU), errors.As(err, &invalidPage), errors.As(err, &invalidSize):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.As(err, &outOfRange):
		return status.Error(codes.OutOfRange, outOfRange.Error())
	default:
		return status.Error(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}
}
