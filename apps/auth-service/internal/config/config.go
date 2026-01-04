package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"port"`
	GRPCPort string `mapstructure:"grpc_port"`
	AppEnv   string `mapstructure:"app_env"`
	AppURL   string `mapstructure:"app_url"`
	JWT      JWTConfig
	MongoDB  MongoConfig
	Redis    RedisConfig
	SMTP     SMTPConfig
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
}

type MongoConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type SMTPConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

func LoadConfig() (config Config, err error) {
	// Read from config file
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
		// Config file not found; ignore error if desired
	}

	// Map environment variables to nested keys
	// e.g. SMTP_HOST -> smtp.host
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}
