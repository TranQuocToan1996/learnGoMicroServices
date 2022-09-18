package main

import (
	"brokerservice/model"
	"fmt"
	"log"
	"net/http"
)

func main() {

	cfg := model.GetConfig()
	// cfg := &model.Config{
	// 	Port: "81",
	// }
	if cfg == nil {
		log.Fatal("empty config")
	}

	mux := routes(cfg)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.Port),
		Handler: mux,
	}

	log.Printf("Starting on port %v", cfg.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
