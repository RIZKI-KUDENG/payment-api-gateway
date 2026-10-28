package payment

import (
	"context"
	"errors"
	"fmt"
	"payment-api-gateway/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)



var (
	ErrNotFound = errors.New("payment intent not found")
)


type Repository interface {
	CreateIntent(ctx context.Context, intent *models.PaymentIntents) error
	GetIntentByID(ctx context.Context, tenantID, id string) (*models.PaymentIntents, error)
	ListIntents(ctx context.Context, tenantID string, limit, offset int) ([]*models.PaymentIntents, error)
	UpdateIntentStatus(ctx context.Context, tenantID, id, status string) error
	CreateAttempt(ctx context.Context, attempt *models.PaymentAttempts) error
	GetLatestAttempt(ctx context.Context, tenantID, paymentIntentID string) (*models.PaymentAttempts, error)
}

type pgRepository struct{
	db *pgxpool.Pool
}


func (r *pgRepository) CreateIntent(ctx context.Context,intent *models.PaymentIntents) error{
	query := `
	INSERT INTO payment_intents(
	tenant_id,
	merchant_account_id,
	merchant_reference,
	amount,
	currency,
	status,
	expires_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		intent.TenantID,
		intent.MerchantAccountID,
		intent.MerchantReference,
		intent.Amount,
		intent.Currency,
		intent.Status,
		intent.ExpiresAt,
	).Scan(&intent.ID, &intent.CreatedAt, &intent.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment intent: %w", err)
	}
	return nil
}

func (r *pgRepository) GetIntentByID(ctx context.Context, tenantID, id string) (*models.PaymentIntents, error){

	query := `
	SELECT id,
	tenant_id,
 	merchant_account_id,
 	merchant_reference,
  	amount,
   currency,
   status,
   expires_at,
   created_at,
   updated_at
   FROM payment_intents
   WHERE tenant_id = $1
   AND id = $2
	`
	var intent models.PaymentIntents
	err := r.db.QueryRow(
		ctx,
		query,
		tenantID,
		id,
	).Scan(
		&intent.ID,
		&intent.TenantID,
		&intent.MerchantAccountID,
		&intent.MerchantReference,
		&intent.Amount,
		&intent.Currency,
		&intent.Status,
		&intent.ExpiresAt,
		&intent.CreatedAt,
		&intent.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return &intent, nil
}

func (r *pgRepository) ListIntents(ctx context.Context, tenantID string, limit, offset int) ([]*models.PaymentIntents, error) {

	query := `
	SELECT id,
	tenant_id,
 	merchant_account_id,
 	merchant_reference,
  	amount,
   currency,
   status,
   expires_at,
   created_at,
   updated_at
   FROM payment_intents
   WHERE tenant_id = $1
   ORDER BY created_at DESC
   LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment intents: %w", err)
	}
	defer rows.Close()

	var list []*models.PaymentIntents

	for rows.Next(){
		var intent models.PaymentIntents

		if err := rows.Scan(
			&intent.ID,
			&intent.TenantID,
			&intent.MerchantAccountID,
			&intent.MerchantReference,
			&intent.Amount,
			&intent.Currency,
			&intent.Status,
			&intent.ExpiresAt,
			&intent.CreatedAt,
			&intent.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment intent: %w", err)
		}
		list = append(list, &intent)
	}
	return list, nil
}

func (r *pgRepository) UpdateIntentStatus(ctx context.Context, tenantID, id, status string) error{
	query := `
	UPDATE payment_intents
	SET status = $1, updated_at = CURRENT_TIMESTAMP
	WHERE tenant_id = $2 AND id = $3
	`

	cmdTag, err := r.db.Exec(ctx, query, status, tenantID, id)
	if err != nil{
		return fmt.Errorf("failed to update payment intent status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgRepository) CreateAttempt(ctx context.Context, attempt *models.PaymentAttempts) error {
	query := `
	INSERT INTO payment_attempts(
	tenant_id,
	payment_intent_id,
	provider_connection_id,
	attempt_number,
	provider_payment_id,
	provider_status,
	status,
	error_code,
	error_message
	) VALUES ($1, $2, $3, $4, $5, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		attempt.TenantID,
		attempt.PaymentIntentID,
		attempt.ProviderConnectionID,
		attempt.AttemptNumber,
		attempt.ProviderPaymentID,
		attempt.ProviderStatus,
		attempt.Status,
		attempt.ErrorCode,
		attempt.ErrorMessage,
	).Scan(&attempt.ID, &attempt.CreatedAt, &attempt.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment attempt: %w", err)
	}
	return nil
}

func (r *pgRepository) GetLatestAttempt(ctx context.Context, tenantID, paymentIntentID string) (*models.PaymentAttempts, error){
	query := `
	SELECT
	id,
	tenant_id,
	payment_intent_id,
	provider_connection_id,
	attempt_number,
	provider_payment_id,
	provider_status,
	status,
	error_code,
	error_message,
	created_at,
	updated_at
	FROM payment_attempts
	WHERE tenant_id = $1 AND payment_intent_id = $2
	ORDER BY created_at DESC
	LIMIT 1
	`

	var att models.PaymentAttempts
	err := r.db.QueryRow(
		ctx, query, tenantID, paymentIntentID,
	).Scan(
			&att.ID,
            &att.TenantID,
            &att.PaymentIntentID,
            &att.ProviderConnectionID,
            &att.AttemptNumber,
            &att.ProviderPaymentID,
            &att.ProviderStatus,
            &att.Status,
            &att.ErrorCode,
            &att.ErrorMessage,
            &att.CreatedAt,
            &att.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest payment attempt: %w", err)
	}

	return &att, nil
}
