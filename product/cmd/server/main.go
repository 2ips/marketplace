package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/2ips/marketplace/product/internal/app/usecase"
	"github.com/2ips/marketplace/product/internal/infrastructure/memory"
	trasnsport "github.com/2ips/marketplace/product/internal/presentation/grpc"
	pb "github.com/2ips/marketplace/product/pkg/api"
	"google.golang.org/grpc"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	listen, err := net.Listen("tcp", config.GRPCListenAddr)
	if err != nil {
		log.Fatalf("Failed to start listener on %s: %v", config.GRPCListenAddr, err)
	}

	if err := run(config, listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func run(config config, listen net.Listener) error {
	productsFile, err := os.Open(config.ProductsJsonPath)
	if err != nil {
		return fmt.Errorf("open products json file %s: %w", config.ProductsJsonPath, err)
	}
	defer func(productsFile *os.File) {
		_ = productsFile.Close()
	}(productsFile)

	productRepository, err := memory.NewProductRepository(productsFile)
	if err != nil {
		return fmt.Errorf("create product repository: %w", err)
	}

	productServer := trasnsport.NewProductServer(
		usecase.NewGetProduct(productRepository),
		usecase.NewListProducts(productRepository),
	)

	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productServer)

	return server.Serve(listen)
}
