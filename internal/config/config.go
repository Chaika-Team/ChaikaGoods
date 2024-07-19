package config

import (
	"github.com/go-kit/log"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

// Config is the application configuration structure that is read from the config file.
type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
}

// StorageConfig is the database configuration structure that is read from the config file.
type StorageConfig struct {
	Host     string `json:"host" env-default:""`
	Port     string `json:"port" env-default:""`
	Database string `json:"database" env-default:""`
	User     string `json:"user" env-default:""`
	Password string `json:"password"`
}

var instance *Config
var once sync.Once

// GetConfig reads the application configuration from the default path.
func GetConfig(log log.Logger) *Config {
	return GetConfigWithPath(log, "config.yml")
}

// GetConfigWithPath reads the application configuration from a given path.
func GetConfigWithPath(logger log.Logger, path string) *Config {
	once.Do(func() {
		logger := log.With(logger, "method", "GetConfig")
		logger.Log("message", "Reading configuration from", "path", path)
		instance = &Config{}
		if err := cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Log("message", "Failed to read configuration", "err", err, "help", help)
			//Fatal err, as the application cannot start without a valid configuration
			panic(err)
		} else {
			logger.Log("message", "Configuration read successfully")
		}
	})
	return instance
}
