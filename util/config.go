package util

import (
	"time"

	"github.com/spf13/viper"
)

//Config stores all configurations of the application.
//The values are read from a config file or environment variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBUrl string `mapstructure:"DB_URL"`
	MigrationUrl string `mapstructure:"MIGRATION_URL"`
	Port string `mapstructure:"PORT_ADDRESS"`
	SymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	SecretKey string `mapstructure:"TOKEN_SECRET_KEY"`
	TokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}