package models

import "time"


type PaymentIntents struct {
	ID string `json:"id"`
	TenantID string `json:"tenant_id"`
	MerchantAccountID string `json:"merchant_account_id"`
	ProviderConnectionID string `json:"provider_connection_id"`
	MerchantReference string `json:"merchant_reference"`
	Amount int `json:"amount"`
	Currency string `json:"currency"`
	Status string `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	TimeStamps
}


type PaymentAttempts struct{
	ID string `json:"id"`
	PaymentIntentID string `json:"payment_intent_id"`
	ProviderConnectionID string `json:"provider_connection_id"`
	AttemptNumber int `json:"attempt_number"`
	ProviderPaymentID string `json:"provider_payment_id"`
	ProviderStatus string `json:"provider_status"`
	Status string `json:"status"`
	ErrorCode string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	TimeStamps
}

type WebhookEvents struct {
	ID string `json:"id"`
	ProviderConnectionID	string `json:"provider_connection_id"`
	ProviderEventID	string `json:"provider_event_id"`
	EventFingerprint string `json:"event_fingerprint"`
	EventType string `json:"event_type"`
	PayloadRef string `json:"payload_ref"`
	ProcessingStatus string `json:"processing_status"`
	ReceivedAt time.Time `json:"received_at"`
	ProcessedAt time.Time `json:"processed_at"`
	RetryCount int `json:"retry_count"`
	LastError string `json:"last_error"`
}

type IdempotencyRecords struct {
	TenantID string `json:"tenant_id"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestHash string `json:"request_hash"`
	ResponseRef string `json:"response_ref"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
