package config

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	DBDialect string
	DBDsn     string
}

func LoadConfig(filepath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filepath)

	v.AutomaticEnv()

	v.SetDefault("DB_DIALECT", "sqlite3")
	v.SetDefault("DB_DSN", "file:session_store.db?_foreign_keys=on")

	if err := v.ReadInConfig(); err != nil {
		logrus.Errorf("failed to read config file: %w", err)
	}

	cfg := &Config{
		DBDialect: v.GetString("DB_DIALECT"),
		DBDsn:     v.GetString("DB_DSN"),
	}

	if cfg.DBDialect == "" || cfg.DBDsn == "" {
		return nil, fmt.Errorf("DB_DIALECT or DB_DSN cannot be empty")
	}

	return cfg, nil
}
