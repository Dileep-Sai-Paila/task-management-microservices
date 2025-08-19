package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver               string `mapstructure:"DB_DRIVER"`
	DBSource               string `mapstructure:"DB_SOURCE"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	UserServiceGRPCAddress string `mapstructure:"USER_SERVICE_GRPC_ADDRESS"`
	RedisAddress           string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found, using environment variables")
		} else {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
