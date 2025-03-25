package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-kit/log/level"

	"github.com/go-kit/log"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the application configuration structure that is read from the config file.
type Config struct {
	Log struct {
		Level string `yaml:"level" env-default:"info" env:"LOG_LEVEL"`
	}
	Listen struct {
		Type   string `yaml:"type" env-default:"port" env:"LISTEN_TYPE"`
		BindIP string `yaml:"bind_ip" env-default:"0.0.0.0" env:"BIND_IP"`
		Port   string `yaml:"port" env-default:"8080" env:"PORT"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
}

// StorageConfig is the database configuration structure that is read from the config file.
type StorageConfig struct {
	Host        string `yaml:"host" env-default:"localhost" env:"DB_HOST"`
	Port        string `yaml:"port" env-default:"5432" env:"DB_PORT"`
	Database    string `yaml:"database" env-default:"chaikagoods" env:"DB_NAME"`
	User        string `yaml:"user" env-default:"postgres" env:"DB_USER"`
	Password    string `yaml:"password" env-required:"true" env:"DB_PASSWORD"`
	Schema      string `yaml:"schema" env-default:"public" env:"DB_SCHEMA"`
	MaxAttempts int    `yaml:"max_attempts" env-default:"5" env:"DB_MAX_ATTEMPTS"`
	MaxConns    int32  `yaml:"max_conns" env-default:"10" env:"DB_MAX_CONNS"`
	MinConns    int32  `yaml:"min_conns" env-default:"2" env:"DB_MIN_CONNS"`
	// Дополнительные параметры пула соединений
	MaxConnLifetime       time.Duration `yaml:"max_conn_lifetime" env-default:"1h" env:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime       time.Duration `yaml:"max_conn_idle_time" env-default:"30m" env:"DB_MAX_CONN_IDLE_TIME"`
	HealthCheckPeriod     time.Duration `yaml:"health_check_period" env-default:"1m" env:"DB_HEALTH_CHECK_PERIOD"`
	MaxConnLifetimeJitter time.Duration `yaml:"max_conn_lifetime_jitter" env-default:"0s" env:"DB_MAX_CONN_LIFETIME_JITTER"`
}

// DSN формирует строку подключения для pgx.
func (c *StorageConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.Schema)
}

var instance *Config
var once sync.Once

// GetConfig reads the application configuration from a given path or from environment variables if the file doesn't exist.
func GetConfig(logger log.Logger, path string) *Config {
	once.Do(func() {
		logger := log.With(logger, "method", "GetConfig")
		instance = &Config{}

		// Проверка наличия файла конфигурации
		if _, err := os.Stat(path); os.IsNotExist(err) {
			_ = level.Info(logger).Log("message", "Config file not found, loading from environment variables", "path", path)
			// Загружаем только из переменных окружения
			if err := cleanenv.ReadEnv(instance); err != nil {
				_ = level.Error(logger).Log("message", "Failed to read configuration from environment variables", "err", err)
				panic(err)
			}
		} else {
			_ = logger.Log("message", "Reading configuration from file", "path", path)
			// Загружаем из файла и окружения
			if err := cleanenv.ReadConfig(path, instance); err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				_ = level.Error(logger).Log("message", "Failed to read configuration from file", "err", err, "help", help)
				panic(err)
			}
		}
		_ = level.Info(logger).Log("message", "Configuration loaded successfully")
	})
	return instance
}
