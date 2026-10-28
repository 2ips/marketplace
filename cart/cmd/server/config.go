package main

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type config struct {
	HTTPListenAddr     string        `env:"HTTP_LISTEN_ADDR" envDefault:"localhost:8080"`
	ProductGRPCAddr    string        `env:"PRODUCT_GRPC_LISTEN_ADDR" envDefault:"localhost:50051"`
	ProductGRPCTimeout time.Duration `env:"PRODUCT_GRPC_TIMEOUT" envDefault:"1s"`
}

func loadConfig() (config, error) {
	var cfg config
	return cfg, env.Parse(&cfg)
}
