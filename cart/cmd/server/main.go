package main

import (
	"log"
	"net/http"

	"github.com/2ips/marketplace/cart/internal/app/usecase"
	"github.com/2ips/marketplace/cart/internal/infrastructure/memory"
	"github.com/2ips/marketplace/cart/internal/infrastructure/productgrpc"
	transport "github.com/2ips/marketplace/cart/internal/presentation/http"
	"github.com/2ips/marketplace/cart/pkg/api"
	pb "github.com/2ips/marketplace/product/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	productClient, err := grpc.NewClient(
		config.ProductGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create product client for addr %s: %v", config.ProductGRPCAddr, err)
	}
	defer func(productClient *grpc.ClientConn) {
		_ = productClient.Close()
	}(productClient)

	productProvider := productgrpc.NewProductProvider(pb.NewProductServiceClient(productClient))
	cartRepository := memory.NewCartRepository()

	cartServer := transport.NewCartServer(
		usecase.NewAddItem(cartRepository, productProvider),
		usecase.NewGetCart(cartRepository),
	)

	handler := api.Handler(api.NewStrictHandler(cartServer, nil))
	server := &http.Server{
		Addr:    config.HTTPListenAddr,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
