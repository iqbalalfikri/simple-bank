package util

import (
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Address string `mapstructure:"address"`
}

type DatabaseConfiguration struct {
	Source string `mapstructure:"source"`
	Driver string `mapstructure:"driver"`
}

type Config struct {
	ServerConfig   ServerConfiguration   `mapstructure:"server"`
	DatabaseConfig DatabaseConfiguration `mapstructure:"database"`
}

func LoadConfig(path, build string) (config Config, err error) {
	viper.AddConfigPath(path)
	if build == "dev" {
		viper.SetConfigName("dev")
	} else {
		viper.SetConfigName("test")
	}
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
