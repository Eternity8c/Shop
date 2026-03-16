package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	StoragePath string
}

func MustLoad() *Config {
	viper.SetConfigFile(os.Getenv("CONFIG_PATH"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config: ", err)
	}

	return &Config{
		Port:        viper.GetString("PORT"),
		StoragePath: viper.GetString("STORAGE_PATH"),
	}
}
