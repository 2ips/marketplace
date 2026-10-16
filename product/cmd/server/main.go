package main

import (
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

	productsFile, err := os.Open(config.ProductsJsonPath)
	if err != nil {
		log.Fatalf("Failed to open products json file %s: %v", config.ProductsJsonPath, err)
	}
	defer func(productsFile *os.File) {
		_ = productsFile.Close()
	}(productsFile)

	productRepository, err := memory.NewProductRepository(productsFile)
	if err != nil {
		log.Fatalf("Failed to create product repository: %v", err)
	}

	productServer := trasnsport.NewProductServer(
		usecase.NewGetProduct(productRepository),
		usecase.NewListProducts(productRepository),
	)

	listen, err := net.Listen("tcp", config.GRPCListenAddr)
	if err != nil {
		log.Fatalf("Failed to start listener on %s: %v", config.GRPCListenAddr, err)
	}

	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, productServer)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
