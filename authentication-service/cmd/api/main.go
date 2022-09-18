package main

import (
	"authentication/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	driver = "pgx"
)

func main() {
	log.Println("Starting authenticate service!")

	conn := connectToDB()
	if conn == nil {
		log.Fatal("can't connect to Postgres")
	}

	// cfg := model.GetConfig(conn)
	// if cfg == nil {
	// 	log.Fatal("empty config")
	// }

	cfg := &model.Config{
		DB:            conn,
		Models:        model.New(conn),
		Port:          "83",
		TimeOutSqlSec: 3,
	}

	mux := routes(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: mux,
	}

	log.Printf("Starting on port %v", cfg.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(conectionString string) (*sql.DB, error) {
	db, err := sql.Open(driver, conectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	connectionString := os.Getenv("DSN")
	var count int

	for {
		db, err := openDB(connectionString)
		if err != nil {
			log.Println("DB not ready...", err)
		} else {
			log.Println("Connected")
			return db
		}

		count++
		if count > 10 {
			return nil
		}

		log.Println("Sleep 3 sec...")
		time.Sleep(time.Second * 3)
	}
}
