package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string         `yaml:"env"`
	TokenTTL time.Duration  `yaml:"token_ttl"`
	Database DatabaseConfig `yaml:"database"`
	Grpc     GrpcConfig     `yaml:"grpc"`
	Secret   string         `yaml:"secret"`
}
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type GrpcConfig struct {
	Port    int    `yaml:"port"`
	Timeout string `yaml:"timeout"`
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

	return &cfg
}
