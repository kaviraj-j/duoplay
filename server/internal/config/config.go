package config

import "github.com/spf13/viper"

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// Load function loads the configs from env file and return Config type reference
func Load(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
