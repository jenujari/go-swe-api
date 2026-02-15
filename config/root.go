package config

import (
	"log"
	"os"

	// "github.com/goforj/godump"

	"github.com/spf13/viper"
)

var (
	logger *log.Logger
	config *Config
)

func init() {
	viper.SetConfigName("conf") // Name of the config file (without extension)
	viper.SetConfigType("yaml") // Config file type
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Initialize variables by unmarshaling into the struct
	config = new(Config)
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// config.Secret.UserName = os.Getenv("ZERODHA_USER")
	// config.Secret.Password = os.Getenv("ZERODHA_PASS")
	// config.Secret.Secret = os.Getenv("ZERODHA_SECRET")
	// config.Secret.ApiKey = os.Getenv("ZERODHA_API_KEY")
	// config.Secret.ApiSecret = os.Getenv("ZERODHA_API_SECRET")
	// config.Secret.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	// config.Secret.POSTGRES_DB = os.Getenv("POSTGRES_DB")
	// config.Secret.POSTGRES_USER = os.Getenv("POSTGRES_USER")

	// init log system
	logger = log.Default()
	logger.SetOutput(os.Stdout)

	// godump.Dump(config)
}

func GetLogger() *log.Logger {
	return logger
}

func GetConfig() *Config {
	return config
}
