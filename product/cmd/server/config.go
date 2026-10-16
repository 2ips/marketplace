package main

import "github.com/caarlos0/env/v11"

type config struct {
	GRPCListenAddr   string `env:"PRODUCT_GRPC_LISTEN_ADDR" envDefault:"localhost:50051"`
	ProductsJsonPath string `env:"PRODUCT_JSON_PATH" envDefault:"products.json.example"`
}

func loadConfig() (config, error) {
	var cfg config
	return cfg, env.Parse(&cfg)
}
