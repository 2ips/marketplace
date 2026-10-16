package main

import "github.com/caarlos0/env/v11"

type config struct {
	HTTPListenAddr  string `env:"HTTP_LISTEN_ADDR" envDefault:"localhost:8080"`
	ProductGRPCAddr string `env:"PRODUCT_GRPC_LISTEN_ADDR" envDefault:"localhost:50051"`
}

func loadConfig() (config, error) {
	var cfg config
	return cfg, env.Parse(&cfg)
}
