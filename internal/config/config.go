package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the application configuration structure that is read from the config file.
type Config struct {
	Log struct {
		Level string `yaml:"level" env-default:"info"`
	}
	Listen struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
}

// StorageConfig is the database configuration structure that is read from the config file.
type StorageConfig struct {
	Host        string `yaml:"host" env-default:"localhost"`
	Port        string `yaml:"port" env-default:"5432"`
	Database    string `yaml:"database" env-default:"chaikagoods"`
	User        string `yaml:"user" env-default:"postgres"`
	Password    string `yaml:"password" env-required:"true"`
	MaxAttempts int    `yaml:"max_attempts" env-default:"5"`
	MaxConns    int32  `yaml:"max_conns" env-default:"10"`
	MinConns    int32  `yaml:"min_conns" env-default:"2"`
	// Дополнительные параметры пула соединений
	MaxConnLifetime       time.Duration `yaml:"max_conn_lifetime" env-default:"1h"`
	MaxConnIdleTime       time.Duration `yaml:"max_conn_idle_time" env-default:"30m"`
	HealthCheckPeriod     time.Duration `yaml:"health_check_period" env-default:"1m"`
	MaxConnLifetimeJitter time.Duration `yaml:"max_conn_lifetime_jitter" env-default:"0s"`
}

// DSN формирует строку подключения для pgx.
func (c *StorageConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

var instance *Config
var once sync.Once

// GetConfig reads the application configuration from the default path.
func GetConfig(logger log.Logger) *Config {
	return GetConfigWithPath(logger, "config.yml")
}

// GetConfigWithPath reads the application configuration from a given path.
func GetConfigWithPath(logger log.Logger, path string) *Config {
	once.Do(func() {
		logger := log.With(logger, "method", "GetConfig")
		_ = logger.Log("message", "Reading configuration from", "path", path)
		instance = &Config{}
		if err := cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			_ = logger.Log("message", "Failed to read configuration", "err", err, "help", help)
			// Fatal err, as the application cannot start without a valid configuration
			panic(err)
		} else {
			_ = logger.Log("message", "Configuration read successfully")
		}
	})
	return instance
}
