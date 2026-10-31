package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"payment-api-gateway/internal/credential"
	"payment-api-gateway/internal/provider"
	"payment-api-gateway/internal/provider/midtrans"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)



type pgRegistry struct{
	db *pgxpool.Pool
	credRepo credential.Repository
	masterKey []byte
}


func NewRegistry(db *pgxpool.Pool, credRepo credential.Repository, masterKey []byte) provider.Registry{
	return &pgRegistry{
		db:  db,
		credRepo: credRepo,
		masterKey: masterKey,
	}
}

func (r *pgRegistry) GetProvider(ctx context.Context, connectionID string) (provider.PaymentProvider, error){
	query := `
	SELECT tenant_id, provider, credentials_ref, environment
	FROM provider_connections WHERE id = $1 AND status = 'active'
	`

	var tenantID, providerName, credentialRef, environment string

	err := r.db.QueryRow(
		ctx,
		query,
		connectionID,
	).Scan(&tenantID, &providerName, &credentialRef, &environment)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, fmt.Errorf("provider not found or inactive")
		}
		return nil, fmt.Errorf("provider connection: %w", err)
	}

	encryptedVal, _,err := r.credRepo.GetCredential(ctx, tenantID, credentialRef)

	if err != nil {
		return nil, fmt.Errorf("get credential: %w", err)
	}


	plainBytes, err := credential.Decrypt(encryptedVal, r.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	var credMap map[string]string

	if err:= json.Unmarshal(plainBytes, &credMap); err != nil{
		return nil, fmt.Errorf("unmarshal credential json: %w", err)
	}

	switch providerName {
		case "midtrans":
		return midtrans.NewClient(
			credMap["server_key"],
				credMap["client_key"],
				environment,
		), nil

		default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}

}
