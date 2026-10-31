package idempotency

import (
	"context"
	"errors"
	"fmt"
	"payment-api-gateway/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)



var ErrRecordNotFound = errors.New("record not found")

type Repository interface{
	GetRecord(ctx context.Context, tenantID, key string) (*models.IdempotencyRecords, error)
	CreateRecord(ctx context.Context, record *models.IdempotencyRecords) error
	UpdateRecord(ctx context.Context, tenantID, key, status, response string) error
	DeleteRecord(ctx context.Context, tenantID, key string) error
}


type pgRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetRecord(ctx context.Context, tenantID, key string) (*models.IdempotencyRecords, error) {

	query := `
	SELECT
	tenant_id,
	idempotency_key,
	request_hash,
	response_ref,
	status,
	expires_at,
	created_at
	FROM idempotency_records
	WHERE tenant_id = $1
	AND idempotency_key = $2
	`

	var record models.IdempotencyRecords

	err := r.db.QueryRow(
		ctx,
		query,
		tenantID,
		key,
	).Scan(
		&record.TenantID,
		&record.IdempotencyKey,
		&record.RequestHash,
		&record.ResponseRef,
		&record.Status,
		&record.ExpiresAt,
		&record.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &record, nil
}

func (r *pgRepository) CreateRecord(ctx context.Context, record *models.IdempotencyRecords) error {

	query := `
	INSERT INTO idempotency_records
	(tenant_id, idempotency_key, request_hash, response_ref, status, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(
		ctx,
		query,
		record.TenantID,
		record.IdempotencyKey,
		record.RequestHash,
		record.ResponseRef,
		record.Status,
		record.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create record: %v", err)
	}
	return nil
}

func (r *pgRepository) UpdateRecord(ctx context.Context, tenantID, key, status, response string) error {
	query := `
	UPDATE idempotency_records
	SET status = $3, response_ref = $4
	WHERE tenant_id = $1
	AND idempotency_key = $2
	`
	cmdTag, err := r.db.Exec(
		ctx,
		query,
		tenantID,
		key,
		status,
		response,
	)
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no record found")
	}
	return nil
}

func (r *pgRepository) DeleteRecord(ctx context.Context, tenantID, key string) error {
	query := `DELETE FROM idempotency_records WHERE tenant_id = $1 AND idempotency_key = $2`
	_, err := r.db.Exec(ctx, query, tenantID, key)
	return err
}
