package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

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

	// Try to read config file, but don't fail if it doesn't exist (for Render deployment)
	_ = viper.ReadInConfig()

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	// Handle PORT environment variable (provided by Render and other cloud platforms)
	// Priority: PORT > SERVER_ADDRESS > default
	port := os.Getenv("PORT")
	if port != "" {
		// Validate port is a number
		if _, err := strconv.Atoi(port); err != nil {
			return config, fmt.Errorf("invalid PORT value: %s", port)
		}
		// Use 0.0.0.0 to listen on all interfaces (required for Render)
		config.ServerAddress = fmt.Sprintf("0.0.0.0:%s", port)
	} else if config.ServerAddress == "" {
		// Default to localhost:8080 if neither PORT nor SERVER_ADDRESS is set
		config.ServerAddress = ":8080"
	}

	return config, nil
}
