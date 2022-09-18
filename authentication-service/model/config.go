package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	connOnce sync.Once
	con      *Config

	ErrMustSingleJson = errors.New("body must have only one json value")
)

const (
	// 1 Mb
	maxReadOneMb = 1 << 20
)

type Config struct {
	DB            *sql.DB `json:"-"`
	Models        Models  `json:"-"`
	Port          string  `json:"port"`
	TimeOutSqlSec int     `json:"timeoutSqlSec,omitempty"`
}

func GetConfig(db *sql.DB) *Config {
	if con != nil {
		return con
	}

	// connOnce.Do(func() {
	// 	setConfig(db)
	// })
	setConfig(db)

	return con
}

func setConfig(db *sql.DB) {
	cfg := &Config{}
	fileConfig := os.Getenv("config")
	if fileConfig == "" {
		fileConfig = "configExample.json"
	}

	buf, _ := os.ReadFile(fileConfig)
	if len(buf) == 0 {
		buf, _ = os.ReadFile("../configExample.json")
		if len(buf) == 0 {
			buf, _ = os.ReadFile("../../configExample.json")
			if len(buf) == 0 {
				buf, _ = os.ReadFile("../../../configExample.json")
			}
		}
	}
	fmt.Println("Once executive: len buf is ", len(buf))
	err := json.Unmarshal(buf, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Port == "" {
		cfg.Port = "83"
		cfg.TimeOutSqlSec = 3
	}

	cfg.DB = db
	cfg.Models = New(db)

	con = cfg
}
