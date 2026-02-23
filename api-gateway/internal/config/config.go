package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	APIGatewayPort     string
	AuthServicesAddr   string
	ProductServiceAddr string
}

func MustLoad() *Config {
	viper.SetConfigFile(os.Getenv("CONFIG_PATH"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		APIGatewayPort:     viper.GetString("GATEWAY_PORT"),
		AuthServicesAddr:   viper.GetString("AUTH_SERVICES_ADDR"),
		ProductServiceAddr: viper.GetString("PRODUCT_SERVICE_ADDR"),
	}
}
