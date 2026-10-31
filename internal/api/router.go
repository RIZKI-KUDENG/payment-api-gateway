package api

import (
	"net/http"
	"os"

	"payment-api-gateway/internal/api/response"
	v1 "payment-api-gateway/internal/api/v1"
	"payment-api-gateway/internal/auth"
	"payment-api-gateway/internal/credential"
	"payment-api-gateway/internal/idempotency"
	"payment-api-gateway/internal/payment"
	"payment-api-gateway/internal/provider/registry"

	"github.com/jackc/pgx/v5/pgxpool"
)


func NewRouter(db *pgxpool.Pool, authenticator auth.APIKeyAuthenticator) http.Handler {

	idempRepo := idempotency.NewRepository(db)
	idempService := idempotency.NewService(idempRepo)
	credRepo := credential.NewRepository(db)
	masterKey := []byte(os.Getenv("APP_ENCRYPTION_KEY"))
	providerRegistry := registry.NewRegistry(db, credRepo, masterKey)

	paymentRepo := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepo, providerRegistry, idempService)

	paymentHandler := v1.NewPaymentHandler(paymentService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteSuccess(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	authMiddleware := auth.TenantAuthMiddleware(authenticator)
	v1Mux := http.NewServeMux()
	v1Mux.HandleFunc("POST /v1/payment-intents", paymentHandler.CreatePaymentIntent)
	v1Mux.HandleFunc("GET /v1/payment-intents/{id}", paymentHandler.GetPaymentIntent)

	mux.Handle("/v1/", authMiddleware(v1Mux))

	return mux
}
