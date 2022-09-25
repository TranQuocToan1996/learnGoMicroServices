package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// Write log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged success",
	}

	return res, nil
}

func (c *Config) gRPCListen() {
	listen, err := net.Listen(tcp, fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("fail to listen gRPC: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{
		Models: c.Models,
	})

	log.Printf("gRPC server starting on port %s", gRpcPort)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("fail to listen gRPC: %v", err)
	}

}
