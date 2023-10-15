package main

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Server struct {
		GrpcPort string `env:"GRPC_PORT" env-default:"8080"`
	}

	Cache struct {
		Capacity int `env:"CACHE_CAP" env-default:"1000"`
	}
}

func NewConfig() (*Config, error) {
	var c Config
	if err := cleanenv.ReadEnv(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
