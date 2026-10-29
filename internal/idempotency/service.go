package idempotency

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"payment-api-gateway/internal/models"
)



var (
	ErrPayloadMissmatch = errors.New("idempotency reused with different payload")
	ErrInProgress = errors.New("idempotent request is still processing")
	ErrConflict = errors.New("idempotent key already claimed")
)



type Service struct{
	repo Repository
}


func NewService(repo Repository) *Service{
	return &Service{repo: repo}
}


func HashPayload(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}


func (s *Service) Begin(
	ctx context.Context,
	tenantID string,
	key string,
	payload []byte,
)(*models.IdempotencyRecords, bool, error) {
	if key == ""{
		return nil, false, errors.New("idempotency key is required")
	}

	requestHash := HashPayload(payload)


	record, err := s.repo.GetRecord(ctx, tenantID, key)

	if err == nil {
		if record.RequestHash != requestHash{
			return nil, false,	ErrPayloadMissmatch
		}
		switch record.Status {
		case "completed":
			return record, true, nil
		case "processing":
			return nil, false, ErrInProgress
		default:
			return nil, false, fmt.Errorf("unknown status : %s", record.Status)
		}
	}

	if !errors.Is(err, ErrRecordNotFound){
		return nil, false, fmt.Errorf("get idempotency record: %w", err)
	}

	record = &models.IdempotencyRecords{
		TenantID: tenantID,
		IdempotencyKey: key,
		RequestHash: requestHash,
		Status: "processing",
	}

	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return nil, false, fmt.Errorf("claim idempotency key: %w", err)
	}

	return record, false, nil
}


func (s *Service) Complete(
	ctx context.Context,
	tenantID string,
	key string,
	responseRef string,
)error {
	if err := s.repo.UpdateRecord(
		ctx,
		tenantID,
		key,
		"completed",
		responseRef,
	); err != nil {
		return fmt.Errorf("complete idempotency record: %w", err)
	}
	return nil
}
