//go:build api

package internal

import (
	"testing"

	"github.com/2ips/marketplace/cart/pkg/api"
	pb "github.com/2ips/marketplace/product/pkg/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCart(t testing.TB) *api.ClientWithResponses {
	t.Helper()

	cfg, err := loadConfig()
	require.NoError(t, err)

	client, err := api.NewClientWithResponses("http://" + cfg.CartHTTPAddr)
	require.NoError(t, err)

	return client
}

func NewProduct(t testing.TB) pb.ProductServiceClient {
	t.Helper()

	cfg, err := loadConfig()
	require.NoError(t, err)

	conn, err := grpc.NewClient(cfg.ProductGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	return pb.NewProductServiceClient(conn)
}
