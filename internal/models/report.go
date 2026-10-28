package models

import (
	"encoding/json"
	"time"
)

type SettlementBatches struct {
	ID                   string    `json:"id"`
	ProviderConnectionID string    `json:"provider_connection_id"`
	ProviderBatchID      string    `json:"provider_batch_id"`
	SettlementDate       time.Time `json:"settlement_date"`
	Currency             string    `json:"currency"`
	GrossAmount          int64     `json:"gross_amount"`
	Fee                  int64     `json:"fee"`
	NetAmount            int64     `json:"net_amount"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
}

type SettlementItems struct {
	ID                    string    `json:"id"`
	SettlementBatchID     string    `json:"settlement_batch_id"`
	ProviderTransactionID string    `json:"provider_transaction_id,omitempty"`
	ProviderPaymentID     string    `json:"provider_payment_id,omitempty"`
	MatchedPaymentID      string    `json:"matched_payment_id,omitempty"`
	GrossAmount           int64     `json:"gross_amount"`
	Fee                   int64     `json:"fee"`
	NetAmount             int64     `json:"net_amount"`
	Currency              string    `json:"currency"`
	MatchStatus           string    `json:"match_status"`
	CreatedAt             time.Time `json:"created_at"`
}

type AuditLogs struct {
	ID           string          `json:"id"`
	TenantID     string          `json:"tenant_id"`
	ActorID      string          `json:"actor_id,omitempty"`
	ActorType    string          `json:"actor_type"`
	Source       string          `json:"source"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   string          `json:"resource_id,omitempty"`
	RequestID    string          `json:"request_id,omitempty"`
	Metadata     json.RawMessage `json:"metadata"`
	CreatedAt    time.Time       `json:"created_at"`
}
