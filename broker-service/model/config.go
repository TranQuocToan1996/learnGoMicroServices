package model

import (
	"brokerservice/event"
	"brokerservice/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	Port     string           `json:"port"`
	RabbitMQ *amqp.Connection `json:"-"`
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
	case logAction:
		// c.log(w, req.Log) old log
		// c.logMessageToRabbitMQ(w, req.Log)
		c.logViaRPC(w, req.Log)
	case mail:
		c.sendMail(w, req.Mail)
	default:
		c.WriteErrJson(w, ErrUnknownAction, http.StatusBadRequest)
	}

}

func (c *Config) sendMail(w http.ResponseWriter, l MailPayload) {
	buf, err := json.Marshal(l)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	newReq, err := http.NewRequest(http.MethodPost, MailServiceURL, bytes.NewBuffer(buf))
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	newReq.Header.Set(ContentType, ApplicationJson)
	client := http.Client{
		Timeout: time.Second * 180,
	}

	res, err := client.Do(newReq)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		c.WriteErrJson(w, ErrUnAuth)
		return
	}

	var payload Response
	payload.Message = fmt.Sprintf("Mail to: %v", l.To)
	c.WriteJson(w, http.StatusOK, payload)

}

// old log func, sending req to log service
func (c *Config) log(w http.ResponseWriter, log LogPayload) {
	buf, err := json.Marshal(&log)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	newReq, err := http.NewRequest(http.MethodPost, LogServiceURL, bytes.NewBuffer(buf))
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	newReq.Header.Set(ContentType, ApplicationJson)
	client := http.Client{
		Timeout: time.Second * 180,
	}

	res, err := client.Do(newReq)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		c.WriteErrJson(w, ErrUnAuth)
		return
	}

	var payload Response
	payload.Message = "logged"
	c.WriteJson(w, http.StatusOK, payload)
}

func (c *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(c.RabbitMQ)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = emitter.Push(context.Background(), string(buf), InfoLog)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) logMessageToRabbitMQ(w http.ResponseWriter, logData LogPayload) {
	err := c.pushToQueue(logData.Name, logData.Data)
	if err != nil {
		log.Println(err)
		c.WriteErrJson(w, err)
		return
	}

	var payload Response
	payload.Message = fmt.Sprintf("logged via RabbitMQ with [name, data]: [%v, %v]", logData.Name, logData.Data)
	c.WriteJson(w, http.StatusOK, &payload)
}

func (c *Config) logViaRPC(w http.ResponseWriter, logData LogPayload) {
	// Connect to RPC server in the logger service through port 5001
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}
	rpcPayload := &RPCPayload{
		Name: logData.Name,
		Data: logData.Data,
	}

	var result string

	// RPCServer is name of struct in server end (logger service) with the method LogInfo
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	payload := Response{
		Message: result,
	}

	c.WriteJson(w, http.StatusOK, &payload)
}

func (c *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var reqPayload Request

	err := c.ReadJson(w, r, &reqPayload)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	defer conn.Close()

	client := logs.NewLogServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: reqPayload.Log.Name,
			Data: reqPayload.Log.Data,
		},
	})
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	var payload Response
	payload.Message = "logged via gRPC"

	c.WriteJson(w, http.StatusOK, &payload)

}

func (c *Config) authenticate(w http.ResponseWriter, authReq AuthPayload) {
	// Debug
	// const (
	// 	pre     = `http://`
	// 	baseUrl = `authentication-service` // name in docker-compose
	// 	// baseUrl      = `localhost:8081` // debug
	// 	route = `/authenticate`
	// 	// fullEndpoint = pre + baseUrl + route // Hover to check
	// )

	// Call to auth service
	buf, err := json.Marshal(authReq)
	if err != nil {
		c.WriteErrJson(w, err)
		return
	}

	// Create req, but not yet send
	newReq, err := http.NewRequest(http.MethodPost, AuthenServiceURL, bytes.NewBuffer(buf))
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
	// var netTransport = &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout: 5 * time.Second,
	// 	}).Dial, // Consider add context
	// 	TLSHandshakeTimeout: 5 * time.Second,
	// }

	// Create client
	client := &http.Client{
		// Transport: nil, // Default Transport
		// Transport: netTransport,
		Timeout: time.Second * 10,
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
	} else if response.StatusCode != http.StatusOK {
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

	err = c.WriteJson(w, http.StatusOK, jsonResp)
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

	w.Header().Set(ContentType, ApplicationJson)
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
