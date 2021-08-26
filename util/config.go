package util

import "github.com/spf13/viper"

// config stores all configurations for the service
// the valuses are ready by viper from a config file or environment variable
type Config struct {
	DBDriver 		string	`mapstructure:"DB_DRIVER"`
	DBSource 		string	`mapstructure:"DB_SOURCE"`
	ServerAddress 	string	`mapstructure:"SERVER_ADDRESS"`
}

// reads configuration from a file if path exists or override with env var 
func LoadConfig(path string) (config Config, err error){
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
