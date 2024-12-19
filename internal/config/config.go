package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer contains the HTTP server configuration.
type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// type MongoDBConfig struct {
// 	URI      string
// 	Database string
// }

// Config represents the application configuration.
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"` // Corrected the YAML tag for HTTPServer
}

// MustLoad loads the configuration from the file and returns it.
func MustLoad() *Config {
	var configPath string

	// Check if the config path is set through the environment variable
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// If not, read it from the command line flag
		flags := flag.String("config", "", "path to configuration file")
		flag.Parse()
		configPath = *flags
		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist %s", configPath)
	}

	var cfg Config

	// Read the config file and load the values into the cfg variable
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	// Return the loaded configuration
	return &cfg
}
