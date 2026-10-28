package main

import (
	"context"
	"log"
	"net/http"
	"payment-api-gateway/internal/app"
	"time"

	"github.com/joho/godotenv"
)



func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	ctx, cancel  := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()
	pool, err := app.NewPool(ctx)
	if err != nil{
		log.Fatal("Failed to connect to database: %w", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadHeaderTimeout: 5 * time.Second ,
	}
	log.Println("Server berjalan di localhost 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: %w", err)
	}
}
