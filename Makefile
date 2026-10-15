GO_MODULE_ROOT := github.com/2ips/marketplace

.PHONY: generate
generate:
	protoc -I . \
		--go_out . --go_opt module=$(GO_MODULE_ROOT) \
		--go-grpc_out . --go-grpc_opt module=$(GO_MODULE_ROOT) \
		product/api/product.proto product/api/page.proto
