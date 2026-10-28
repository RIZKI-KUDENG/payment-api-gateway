package models

import "time"

type SettlementBatches struct {
	ID                   string    `json:"id"`
	ProviderConnectionID string    `json:"provider_connection_id"`
	ProviderBatchID      string    `json:"provider_batch_id"`
	SettlementDate       time.Time `json:"settlement_date"`
	Currency             string    `json:"currency"`
	GrossAmount          int       `json:"gross_amount"`
	Fee                  int       `json:"fee"`
	NetAmount            int       `json:"net_amount"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
}

type SettlementItems struct {
	ID                    string    `json:"id"`
	SettlementBatchID     string    `json:"settlement_batch_id"`
	ProviderTransactionID string    `json:"provider_transaction_id"`
	ProviderPaymentID     string    `json:"provider_payment_id"`
	MatchedPaymentID      string    `json:"matched_payment_id"`
	GrossAmount           int       `json:"gross_amount"`
	Fee                   int       `json:"fee"`
	NetAmount             int       `json:"net_amount"`
	Currency              string    `json:"currency"`
	MatchStatus           string    `json:"match_status"`
	CreatedAt             time.Time `json:"created_at"`
}

type AuditLogs struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	ActorID      string    `json:"actor_id"`
	ActorType    string    `json:"actor_type"`
	Source       string    `json:"source"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	RequestID    string    `json:"request_id"`
	Metadata     string    `json:"metadata"`
	CreatedAt    time.Time `json:"created_at"`
}
