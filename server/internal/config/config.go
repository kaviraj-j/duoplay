package config

import "github.com/spf13/viper"

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	JwtSecret     string `mapstructure:"JWT_SECRET"`
}

// Load function loads the configs from env file and return Config
func Load(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}
	err := viper.Unmarshal(&config)
	return config, err
}
