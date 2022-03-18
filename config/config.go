package config

import (
	"fmt"

	"github.com/sarmerer/go-crypto-dashboard/model"
	"github.com/spf13/viper"
)

var (
	Portfolios []*model.Portfolio

	DefaultDBPath string
)

const (
	DefaultDBName string = "dashboard.db"
)

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	DefaultDBPath = viper.GetString("db_path")

	if err := viper.UnmarshalKey("portfolios", &Portfolios); err != nil {
		return fmt.Errorf("failed to unmarshal portfolios: %v", err)
	}

	return nil
}
