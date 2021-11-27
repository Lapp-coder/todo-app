package config

import (
	"net/http"
	"os"

	"github.com/spf13/viper"
)

const configName = "config"

type Config struct {
	Server
	Service
	PostgresDB
}

type Server struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Handler        http.Handler
	MaxHeaderBytes int `mapstructure:"max_header_bytes"`
	ReadTimeout    int `mapstructure:"read_timeout"`
	WriteTimeout   int `mapstructure:"write_timeout"`
}

type Service struct {
	TokenTTL   int `mapstructure:"token_ttl"`
	SigningKey string
	Salt       string
}

type PostgresDB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func New(configPath string) (Config, error) {
	return parseConfig(configPath)
}

func parseConfig(configPath string) (Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	if err := loadEnv(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("server", &cfg.Server); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("service", &cfg.Service); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgres_db", &cfg.PostgresDB); err != nil {
		return err
	}

	return nil
}

func loadEnv(cfg *Config) error {
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		return errPostgresPasswordIsEmpty
	}

	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" {
		return errSigningKeyIsEmpty
	}

	salt := os.Getenv("SALT")
	if salt == "" {
		return errSaltIsEmpty
	}

	cfg.PostgresDB.Password = postgresPassword
	cfg.Service.SigningKey = signingKey
	cfg.Salt = salt

	return nil
}
