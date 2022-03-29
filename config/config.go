package config

import (
	"fmt"

	"github.com/sarmerer/go-crypto-dashboard/driver/sqlite3"
	"github.com/sarmerer/go-crypto-dashboard/model"
	"github.com/spf13/viper"
)

const (
	DefaultAPIPort int = 3000
)

var (
	DBDriver      string
	SQLite3Config sqlite3.Config
	APIPort       int = DefaultAPIPort
)

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	fields := map[string]interface{}{
		"driver":        &DBDriver,
		"driver_config": &SQLite3Config,
		"api_port":      &APIPort,
	}

	for field, ptr := range fields {
		if err := viper.UnmarshalKey(field, ptr); err != nil {
			return fmt.Errorf("failed to unmarshal key %s: %v", field, err)
		}
	}

	return nil
}

func GetPortfolios() ([]*model.Portfolio, error) {
	portfolios := []*model.Portfolio{}

	if err := viper.UnmarshalKey("portfolios", &portfolios); err != nil {
		return nil, fmt.Errorf("failed to unmarshal portfolios: %v", err)
	}

	return portfolios, nil
}
