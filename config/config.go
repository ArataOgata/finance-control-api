package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env          string     `mapstructure:"env" env-required:"true"`
	Database_URL string     `mapstructure:"database_URL" env-required:"true"`
	HttpServer   HttpServer `mapstructure:"http_server"`
}

type HttpServer struct {
	Address     string        `mapstructure:"address" default:"localhost:8080"`
	Timeout     time.Duration `mapstructure:"timeout" default:"4s"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout" default:"60"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	viper.SetConfigName("local.yml")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/Users/arataogata/dev/go/lesson1/go-api/")

	viper.SetDefault("http_server.address", "localhost:8080")
	viper.SetDefault("http_server.timeout", "4s")
	viper.SetDefault("http_server.idle_timeout", "60s")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using ENV only: %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
