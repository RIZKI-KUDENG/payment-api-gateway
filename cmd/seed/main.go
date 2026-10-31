package main

import (
	"context"
	"encoding/hex"
	"log"
	"os"
	"time"

	"payment-api-gateway/internal/app"
	"payment-api-gateway/internal/credential"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not loaded, using system env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := app.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	masterKey := os.Getenv("APP_ENCRYPTION_KEY")
	key, err := hex.DecodeString(masterKey)
	if err != nil {
		log.Fatalf("Invalid APP_ENCRYPTION_KEY: %v", err)
	}
	if len(key) != 32 {
		log.Fatalf("APP_ENCRYPTION_KEY must be exactly 32 bytes, got %d", len(masterKey))
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	clientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	merchantID := os.Getenv("MIDTRANS_MERCHANT_ID")

	if serverKey == "" || clientKey == "" {
		log.Fatal("MIDTRANS_SERVER_KEY and MIDTRANS_CLIENT_KEY must be set in .env")
	}
	if merchantID == "" {
		merchantID = "G906365427"
	}

	credRepo := credential.NewRepository(pool)
	credService := credential.NewService(credRepo, key)

	connectionID, err := credService.StoreMidtransCredential(
		ctx,
		"a0000000-0000-0000-0000-000000000001", // tenantID
		"b0000000-0000-0000-0000-000000000001", // merchantAccountID
		serverKey,
		clientKey,
		merchantID,
	)
	if err != nil {
		log.Fatalf("Gagal simpan credential: %v", err)
	}

	log.Printf("Berhasil! Credential Midtrans tersimpan rapi dan terenkripsi. Connection ID: %s\n", connectionID)
}
