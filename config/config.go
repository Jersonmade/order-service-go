package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`

	Kafka struct {
		Brokers   []string `mapstructure:"brokers"`
		Topic     string   `mapstructure:"topic"`
		GroupID   string   `mapstructure:"group_id"`
		Partition int      `mapstructure:"partition"`
		MinBytes  int      `mapstructure:"min_bytes"`
		MaxBytes  int      `mapstructure:"max_bytes"`
	} `mapstructure:"kafka"`

	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout_seconds"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.order-service-go")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	return &cfg
}
