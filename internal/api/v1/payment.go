package v1

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"payment-api-gateway/internal/api/response"
	"payment-api-gateway/internal/auth"
	"payment-api-gateway/internal/idempotency"
	"payment-api-gateway/internal/payment"
	"time"
)


type PaymentHandler struct {
	paymentService *payment.Service
}

func NewPaymentHandler(paymentService *payment.Service) *PaymentHandler{
	return &PaymentHandler{paymentService: paymentService}
}


type CreatePaymentIntentRequest struct{
	MerchantAccountID string `json:"merchant_account_id"`
	MerchantReference string `json:"merchant_reference"`
	Amount int64 `json:"amount"`
	Currency string `json:"currency"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (h *PaymentHandler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request){
	tenantID, err := auth.GetTenantID(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized tenant")
		return
	}

	idempKey := r.Header.Get("Idempotency-Key")
	if idempKey == ""{
		response.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Header 'Idempotency-Key is required")
		return
	}


	payloadBytes, err := io.ReadAll(r.Body)

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Failed to read request body")
		return
	}

	var req CreatePaymentIntentRequest

	if err := json.Unmarshal(payloadBytes, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return
	}

	if req.Amount <= 0 {
		response.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Amount must be greater than 0")
		return
	}

	if req.Currency == ""{
		req.Currency = "IDR"
	}

	intent, err := h.paymentService.CreatePaymentIntent(
		r.Context(),
		payment.CreateIntentInput{
			TenantID: tenantID.String(),
			MerchantAccountID: req.MerchantAccountID,
			MerchantReference: req.MerchantReference,
			Amount: req.Amount,
			Currency: req.Currency,
			ExpiresAt: req.ExpiresAt,
		},
		payment.IdempotencyInput{
			TenantID: tenantID.String(),
			Key: idempKey,
			Payload: payloadBytes,
		},
	)

	if err != nil {
		if errors.Is(err, idempotency.ErrPayloadMissmatch){
			response.WriteError(w, http.StatusConflict, "IDEMPOTENCY_CONFLICT", "Idempotency key reused with different payload")
			return
		}

		if errors.Is(err, idempotency.ErrInProgress){
			response.WriteError(w, http.StatusConflict, "IDEMPOTENCY_IN_PROGRESS", "Request with this idempotency key is still processing")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusCreated, intent)
}


func (h *PaymentHandler) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	tenantID, err := auth.GetTenantID(r.Context())

	if err != nil{
		response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Failed to get tenant ID")
		return
	}

	intentID := r.PathValue("id")
	if intentID == ""{
		response.WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "Payment intent ID is_required")
		return
	}

	intent, err := h.paymentService.GetPaymentIntent(
		r.Context(),
		tenantID.String(),
		intentID,
	)

	if err != nil {
		if errors.Is(err, payment.ErrNotFound){
			response.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Payment intent not found")
			return
		}

		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, intent)
}
