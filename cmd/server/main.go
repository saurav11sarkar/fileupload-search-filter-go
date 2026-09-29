package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/saurav11sarkar/001practic/internal/app"
	"github.com/saurav11sarkar/001practic/internal/config"
	databse "github.com/saurav11sarkar/001practic/internal/database"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}

	db, err := databse.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler, err := app.NewHandler(db, cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Hello, World!")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
