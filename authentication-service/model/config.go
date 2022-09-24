package model

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
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

var (
	ErrInvalidCredential = errors.New("invalid credential")
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

func (c *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	req := &ReqAuth{}

	err := c.ReadJson(w, r, req)
	if err != nil {
		c.WriteErrJson(w, err, http.StatusBadRequest)
		return
	}

	// LOL for debug in dev only
	if req.Email == "admin@example.com" && req.Password == "verysecret" {
		c.WriteJson(w, http.StatusOK, &Response{
			Error:   false,
			Message: fmt.Sprintf("Testing Logged in user with email testing %v", req.Email),
		})
		return
	}

	user, err := c.Models.User.GetByEmail(context.Background(), req.Email)
	if err != nil {
		c.WriteErrJson(w, err, http.StatusBadRequest)
		return
	}

	match, err := user.PasswordMatches(req.Password)
	if err != nil || !match {
		c.WriteErrJson(w, err, http.StatusBadRequest)
		return
	}

	// log
	err = c.logRequest("authentication", fmt.Sprintf("%s loggin", user.Email))
	if err != nil || !match {
		c.WriteErrJson(w, err, http.StatusBadRequest)
		return
	}
	resp := &Response{
		Error:   false,
		Message: fmt.Sprintf("Logged in user with email %v", user.Email),
		Data:    user,
	}

	c.WriteJson(w, http.StatusOK, resp)

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

func (c *Config) logRequest(name, data string) error {

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	buf, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	logServiceURL := "http://logger-service/log"

	newReq, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 180 * time.Second,
	}

	_, err = client.Do(newReq)
	if err != nil {
		return err
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
