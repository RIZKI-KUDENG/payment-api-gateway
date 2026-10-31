package credential

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Repository interface {
    SaveCredential(
        ctx context.Context,
        tenantID string,
        encryptedVal []byte,
        keyVersion int,
    ) (string, error)

    GetCredential(
        ctx context.Context,
        tenantID string,
        id string,
    ) ([]byte, int, error)

    CreateProviderConnection(
        ctx context.Context,
        conn ProviderConnectionInput,
    ) (string, error)
}


type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) SaveCredential(
    ctx context.Context,
    tenantID string,
    encryptedVal []byte,
    keyVersion int,
) (string, error) {

    query := `
        INSERT INTO provider_credentials (
            tenant_id,
            encrypted_value,
            encryption_key_version
        )
        VALUES ($1, $2, $3)
        RETURNING id
    `

    var credentialID string

    err := r.db.QueryRow(
        ctx,
        query,
        tenantID,
        encryptedVal,
        keyVersion,
    ).Scan(&credentialID)

    if err != nil {
        return "", fmt.Errorf("save credential: %w", err)
    }

    return credentialID, nil
}
type ProviderConnectionInput struct{
	TenantID string
	MerchantAccountID string
	Provider string
	ProviderAccountID string
	Name string
	CredentialsRef string
	Environment string
}

func (r *pgRepository) CreateProviderConnection(ctx context.Context, input ProviderConnectionInput) (string, error) {
	queryConn := `
		INSERT INTO provider_connections (
		tenant_id, provider, provider_account_id, name, credentials_ref, environment,
		status
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'active')
		RETURNING id
		`
		var connectionID string

		err := r.db.QueryRow(
			ctx,
			queryConn,
			input.TenantID,
			input.Provider,
			input.ProviderAccountID,
			input.Name,
			input.CredentialsRef,
			input.Environment,
		).Scan(&connectionID)

		if err != nil {
			return "", fmt.Errorf("create provider connection: %w", err)
		}

		queryMPC := `
			INSERT INTO merchant_provider_connections (tenant_id, merchant_account_id,provider_connection_id, status, is_default)
			VALUES($1, $2, $3,'active', true)
			ON CONFLICT (tenant_id, merchant_account_id, provider_connection_id) DO UPDATE SET is_default = true,
			updated_at = CURRENT_TIMESTAMP
		`

		if _, err := r.db.Exec(ctx, queryMPC, input.TenantID, input.MerchantAccountID, connectionID); err != nil {
			return "", fmt.Errorf("create provider connection: %w", err)
		}

		return connectionID, nil
}

func (r *pgRepository) GetCredential(
    ctx context.Context,
    tenantID string,
    id string,
) ([]byte, int, error) {

    query := `
        SELECT encrypted_value, encryption_key_version
        FROM provider_credentials
        WHERE tenant_id = $1
          AND id = $2
    `

    var encryptedValue []byte
    var keyVersion int

    err := r.db.QueryRow(ctx, query, tenantID, id).Scan(
        &encryptedValue,
        &keyVersion,
    )

    if err != nil {
        return nil, 0, fmt.Errorf("get credential: %w", err)
    }

    return encryptedValue, keyVersion, nil
}
