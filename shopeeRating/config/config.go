package config

import (
	"encoding/json"
	"os"
	"sync"
)

var configOnce sync.Once = sync.Once{}

type (
	// Main config for dependency injection
	Config struct {
		App  `json:"app"`
		HTTP `json:"http"`
		Log  `json:"logger"`
		API  `json:"api"`
		// PG   `json:"postgres"`
	}

	App struct {
		Name    string `json:"name"    env:"APP_NAME"`
		Version string `json:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `json:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `json:"log_level"   env:"LOG_LEVEL"`
	}

	API struct {
		Limit   int    `json:"limit"   env:"LIMIT"`
		BaseURL string `json:"baseURL"   env:"BASE_URL"`
	}

	// // Poestgres
	// PG struct {
	// 	PoolMax int    `env-required:"true" json:"pool_max" env:"PG_POOL_MAX"`
	// 	URL     string `env-required:"true"                 env:"PG_URL"`
	// }
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	var err error
	var buf []byte
	configOnce.Do(func() {
		buf, err = os.ReadFile("./config/config.json")
		if err != nil {
			buf, err = os.ReadFile("./config.json")
			if err != nil {
				buf, err = os.ReadFile("../config/config.json")
				if err != nil {
					buf, err = os.ReadFile("../../config/config.json")
				}
			}
		}
		err = json.Unmarshal(buf, cfg)
	})

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
