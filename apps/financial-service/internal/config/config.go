package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	AppEnv         string `mapstructure:"APP_ENV"`
	MongoURI       string `mapstructure:"MONGODB_URI"`
	MongoDatabase  string `mapstructure:"MONGODB_DATABASE"`
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      string `mapstructure:"REDIS_PORT"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	JWTSecret           string `mapstructure:"JWT_SECRET"`
	AuthServiceURL      string `mapstructure:"AUTH_SERVICE_URL"`
	AuthServiceGRPCAddr string `mapstructure:"AUTH_SERVICE_GRPC_ADDR"`
}

func LoadConfig() (config Config, err error) {
	viper.SetDefault("PORT", "8082")
	viper.SetDefault("APP_ENV", "production")
	viper.SetDefault("MONGODB_URI", "mongodb://mongodb:27017")
	viper.SetDefault("MONGODB_DATABASE", "gspend")
	viper.SetDefault("REDIS_HOST", "redis")
	viper.SetDefault("REDIS_PORT", "6379")

	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}
