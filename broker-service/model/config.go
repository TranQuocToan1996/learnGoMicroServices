package model

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	Port string `json:"port"`
}

func (c *Config) Broker(w http.ResponseWriter, r *http.Request) {
	req := &Response{
		Message: "Hit the message broker",
	}

	c.WriteJson(w, http.StatusOK, req)
}

func (c *Config) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {

	r.Body = http.MaxBytesReader(w, r.Body, maxReadOneMb)

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	// Decode to noname struct to check whether exist other json datas
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return ErrMustSingleJson
	}

	return nil
}

func (c *Config) WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// TODO: Need handle other index of headers
	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set(ContentType, appicationJsonMIME)
	w.WriteHeader(status)

	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) WriteErrJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := &Response{
		Error:   true,
		Message: err.Error(),
	}

	return c.WriteJson(w, statusCode, payload)
}

func GetConfig() *Config {
	if con != nil {
		return con
	}

	connOnce.Do(func() {
		readConfig()
	})

	return con
}

func readConfig() {
	cfg := &Config{}
	// fileConfig := os.Getenv("config")
	// if fileConfig == "" {
	// 	fileConfig = "configExample.json"
	// }

	// buf, _ := os.ReadFile(fileConfig)
	// if len(buf) == 0 {
	// 	buf, _ = os.ReadFile("../configExample.json")
	// 	if len(buf) == 0 {
	// 		buf, _ = os.ReadFile("../../configExample.json")
	// 		if len(buf) == 0 {
	// 			buf, _ = os.ReadFile("../../../configExample.json")
	// 		}
	// 	}
	// }
	// fmt.Println("len buf", len(buf))
	// err := json.Unmarshal(buf, cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if cfg.Port == "" {
		cfg.Port = "81"
	}

	con = cfg
}
