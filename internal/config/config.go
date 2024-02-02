package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string         `yaml:"env" env:"env"`
	TokenTTL time.Duration  `yaml:"token_ttl" env:"token_ttl"`
	Database DatabaseConfig `yaml:"database"`
	Grpc     GrpcConfig     `yaml:"grpc"`
	Secret   string         `yaml:"secret" env:"secret"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"database_host"`
	Port     int    `yaml:"port" env:"database_port"`
	Username string `yaml:"username" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Name     string `yaml:"name" env:"POSTGRES_DB"`
}

type GrpcConfig struct {
	Port    int    `yaml:"port" env:"grpc_port"`
	Timeout string `yaml:"timeout" env:"grpc_timeout"`
}

func New() *Config {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	if path == "" {
		panic("path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file is not exists")
	}
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("Failed to read config")
	}
	fmt.Println(cfg)
	return &cfg
}
