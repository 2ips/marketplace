GO_MODULE_ROOT := github.com/2ips/marketplace

.PHONY: generate
generate: generate-product generate-cart

.PHONY: generate-product
generate-product:
	protoc -I . \
		--go_out . --go_opt module=$(GO_MODULE_ROOT) \
		--go-grpc_out . --go-grpc_opt module=$(GO_MODULE_ROOT) \
		product/api/product.proto product/api/page.proto

.PHONY: generate-cart
generate-cart:
	npx -y @redocly/cli@1.34.0 bundle cart/api/cart.yml -o cart/api/cart.gen.yml
	cd cart && go tool oapi-codegen -config oapi-codegen.yml api/cart.gen.yml
	rm -f cart/api/cart.gen.yml

.PHONY: run-all
run-all: run-product run-cart

.PHONY: run-product
run-product:
	cd product && go run ./cmd/server


.PHONY: run-cart
run-cart:
	cd cart && go run ./cmd/server
