package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	connOnce sync.Once
	con      *Config

	ErrMustSingleJson = errors.New("body must have only one json value")
	ErrUnknownAction  = errors.New("invalid kind of action")
	ErrUnAuth         = errors.New("unauthenticate")
	ErrCallingAuth    = errors.New("cant calling to auth service")
)

const (
	maxReadOneMb = 1 << 20 // 1 Mb
)

type Config struct {
	Port string `json:"port"`
}

func (c *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	req := &Request{}

	err := c.ReadJson(w, r, req)
	if err != nil {
		c.WriteErrJson(w, err, http.StatusBadRequest)
		return
	}

	switch req.Action {
	case auth:
		c.authenticate(w, req.Auth)
	default:
		c.WriteErrJson(w, ErrUnknownAction, http.StatusBadRequest)
	}

	c.WriteJson(w, http.StatusOK, req)
}

func (c *Config) authenticate(w http.ResponseWriter, authReq AuthPayload) {

	const (
		pre          = `http://`
		baseUrl      = `authentication-service`
		route        = `/authenticate`
		fullEndpoint = pre + baseUrl + route // Hover to check
	)

	// Call to auth service
	buf, err := json.Marshal(authReq)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	// Create req, but not yet send
	newReq, err := http.NewRequest(http.MethodPost, fullEndpoint, bytes.NewBuffer(buf))
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	/*
			timeout with context
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Error("Request error", err)
		}

		resp, err := http.DefaultClient.Do(req.WithContext(ctx)) */

	// More info: request-response cycle is constituted up of Dialer, TLS Handshake, Request Header, Request Body, Response Header and Response Body timeouts
	// Make customer Transport
	// https://itnext.io/http-request-timeouts-in-go-for-beginners-fe6445137c90
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial, // Consider add context
		TLSHandshakeTimeout: 5 * time.Second,
	}

	// Create client
	client := &http.Client{
		// Transport: nil, // Default Transport
		Transport: netTransport,
		Timeout:   time.Second * 10,
	}

	// Send req
	response, err := client.Do(newReq)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	defer response.Body.Close()

	// Check http status resp
	if response.StatusCode == http.StatusUnauthorized {
		c.WriteErrJson(w, ErrUnAuth)
		return
	} else if response.StatusCode != http.StatusAccepted {
		c.WriteErrJson(w, ErrCallingAuth)
		return
	}

	// read body resp
	jsonResp := &Response{}
	err = json.NewDecoder(response.Body).Decode(jsonResp)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	if jsonResp.Error {
		c.WriteErrJson(w, errors.New(jsonResp.Message), http.StatusUnauthorized)
		return
	}

	err = c.WriteJson(w, http.StatusAccepted, jsonResp)
	if err != nil {
		c.WriteErrJson(w, err)
	}
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

	// connOnce.Do(func() {
	// 	readConfig()
	// })
	readConfig()

	return con
}

func readConfig() {
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
		cfg.Port = "81"
	}

	con = cfg
}
