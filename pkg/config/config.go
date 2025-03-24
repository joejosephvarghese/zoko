package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config struct to hold environment variables
type Config struct {
	Port               string `mapstructure:"PORT" validate:"required"`
	DBHost             string `mapstructure:"DB_HOST" validate:"required"`
	DBName             string `mapstructure:"DB_NAME" validate:"required"`
	DBUser             string `mapstructure:"DB_USER" validate:"required"`
	DBPort             string `mapstructure:"DB_PORT" validate:"required"`
	DBPassword         string `mapstructure:"DB_PASSWORD" validate:"required"`
	JwtKey             string `mapstructure:"JWT_KEY" validate:"required"`
	GoAuthClientId     string `mapstructure:"GOAUTH_CLIENT_ID"`
	GoAuthClientSecret string `mapstructure:"GOAUTH_CLIENT_SECRET"`
}

// List of required environment variables
var envs = []string{
	"PORT", "DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "JWT_KEY",
}

// LoadConfig initializes configuration
func LoadConfig() (Config, error) {
	var config Config

	// Check if .env file exists and print its contents for debugging
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println(".env file does not exist, using system environment variables.")
	} else {
		fmt.Println(".env file found.")
		file, err := os.ReadFile(".env")
		if err == nil {
			fmt.Println("Contents of .env:\n", string(file))
		} else {
			fmt.Println("Failed to read .env file:", err)
		}
	}

	// Set environment variable key replacer
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Load system environment variables
	viper.AutomaticEnv()

	// Read .env file
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: Failed to read .env file, using only system environment variables.")
	}

	// Set default values
	viper.SetDefault("DB_PORT", 5432) // Ensure default DB_PORT is an integer

	// Bind environment variables
	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	// Print all loaded environment variables
	fmt.Println("Loaded Env Variables:", viper.AllSettings())

	// Unmarshal config values
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// Ensure DB_PORT is a stringified integer
	config.DBPort = fmt.Sprintf("%d", viper.GetInt("DB_PORT"))

	// Validate config struct
	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	// Debugging - Print final database connection string (without password for security)
	fmt.Printf("Database DSN: host=%s user=%s dbname=%s port=%s sslmode=disable\n",
		config.DBHost, config.DBUser, config.DBName, config.DBPort)

	return config, nil
}
