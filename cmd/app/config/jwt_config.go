package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	varJWTSecret = "CONFIG_JWT_SECRET"
)

type JWTConfig struct {
	Secret string `end:"CONFIG_JWT_SECRET"`
}

func CreateJWTConfig() (JWTConfig, error) {
	var cfg JWTConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
