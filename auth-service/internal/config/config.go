package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	StoragePath string
	Port        int
	TokenTTL    time.Duration
	Timeout     time.Duration
	Secret      string
}

func MustLoad() *Config {
	viper.SetConfigFile(os.Getenv("CONFIG_PATH"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config:", err)
	}

	return &Config{
		StoragePath: viper.GetString("STORAGE_PATH"),
		Port:        viper.GetInt("port"),
		TokenTTL:    viper.GetDuration("TOKEN_TTL"),
		Timeout:     viper.GetDuration("timeout"),
		Secret:      viper.GetString("secret"),
	}
}
