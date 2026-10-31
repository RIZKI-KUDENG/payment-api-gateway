package models

import "time"

type TimeStamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tenant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	TimeStamps
}

type MerchantAccounts struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	TimeStamps
}

type ProviderConnections struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenant_id"`
	Provider          string `json:"provider"`
	ProviderAccountID string `json:"provider_account_id"`
	Name              string `json:"name"`
	CredentialsRef    string `json:"credentials_ref"`
	Environment       string `json:"environment"`
	Status            string `json:"status"`
	TimeStamps
}

type MerchantProviderConnections struct {
	TenantID             string `json:"tenant_id"`
	MerchantAccountID    string `json:"merchant_account_id"`
	ProviderConnectionID string `json:"provider_connection_id"`
	Status               string `json:"status"`
	IsDefault            bool   `json:"is_default"`
	TimeStamps
}

type ProviderCredentials struct {
	ID string `json:"id"`
	TenantID string `json:"tenant_id"`
	EncryptedValue []byte `json:"encrypted_value"`
	EncryptedKeyVersion int `json:"encrypted_key_version"`
	TimeStamps
}
