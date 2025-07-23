// configs/config.go
package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL    string `mapstructure:"DATABASE_URL"`
	Port           string `mapstructure:"PORT"`
	AutoMigrate    bool   `mapstructure:"AUTO_MIGRATE"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
	JWTExpiryHours int    `mapstructure:"JWT_EXPIRY_HOURS"`
}

func LoadConfig() (config Config, err error) {
	v := viper.New()

	v.AddConfigPath("./configs") // Look for config in the current directory
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv() // Read environment variables that match

	err = v.ReadInConfig()
	if err != nil {
		// Handle the case where .env file is not found, but env vars might be set
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: .env file not found: %v", err)
		}
	}

	// Set default value for AUTO_MIGRATE if not provided
	v.SetDefault("AUTO_MIGRATE", false)
	v.SetDefault("JWT_EXPIRY_HOURS", 24) // Default to 24 hours

	err = v.Unmarshal(&config)
	return
}
