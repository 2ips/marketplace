package main

import (
	"fmt"
	"log"
	"net"
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

	listen, err := net.Listen("tcp", config.HTTPListenAddr)
	if err != nil {
		log.Fatalf("Failed to start listener on %s: %v", config.HTTPListenAddr, err)
	}

	if err := run(config, listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func run(config config, listen net.Listener) error {
	productClient, err := grpc.NewClient(
		config.ProductGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("create product client for addr %s: %w", config.ProductGRPCAddr, err)
	}
	defer func(productClient *grpc.ClientConn) {
		_ = productClient.Close()
	}(productClient)

	productProvider := productgrpc.NewProductProvider(pb.NewProductServiceClient(productClient))
	cartRepository := memory.NewCartRepository()

	cartServer := transport.NewCartServer(
		usecase.NewAddItem(cartRepository, productProvider, config.ProductGRPCTimeout),
		usecase.NewGetCart(cartRepository, productProvider, config.ProductGRPCTimeout),
	)

	swaggerSpec, err := api.GetSpec()
	if err != nil {
		return fmt.Errorf("load swagger spec: %w", err)
	}
	swaggerJSON, err := swaggerSpec.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal swagger spec: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /docs/{$}", serveBytes("text/html; charset=utf-8", docsHTML))
	mux.Handle("GET /swagger.json", serveBytes("application/json", swaggerJSON))

	handler := api.HandlerFromMux(api.NewStrictHandler(cartServer, nil), mux)
	server := &http.Server{Handler: handler}

	return server.Serve(listen)
}

func serveBytes(contentType string, body []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(body)
	})
}
