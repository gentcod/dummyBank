package util

import "github.com/spf13/viper"

//Config stores all configurations of the application.
//The values are read from a config file or environment variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBUrl string `mapstructure:"DB_URL"`
	Port string `mapstructure:"PORT_ADDRESS"`
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