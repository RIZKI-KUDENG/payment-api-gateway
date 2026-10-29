package payment

import (
	"context"
	"errors"
	"fmt"
	"payment-api-gateway/internal/idempotency"
	"payment-api-gateway/internal/models"
	"payment-api-gateway/internal/provider"
	"time"
)

type Service struct{
	repo Repository
	registry provider.Registry
	idempotencyService *idempotency.Service
}


func NewService(repo Repository, registry provider.Registry, idempotencyService *idempotency.Service) *Service{
	return &Service{repo: repo, registry: registry, idempotencyService: idempotencyService}
}

type CreateIntentInput struct {
	TenantID string
	MerchantAccountID string
	MerchantReference string
	Amount int64
	Currency string
	ExpiresAt *time.Time
}
type IdempotencyInput struct {
	TenantID string
	Key string
	Payload []byte
}
type SyncPaymentResult struct{
	PaymentIntent *models.PaymentIntents `json:"payment_intent"`
	LatestAttempt *models.PaymentAttempts `json:"latest_attempt,omitempty"`
	RawProviderStatus string `json:"raw_provider_status"`
}

func (s *Service) CreatePaymentIntent(
	ctx context.Context,
	input CreateIntentInput,
	idempotencyInput IdempotencyInput,
)(*models.PaymentIntents, error){
	 record, err := s.idempotencyService.Begin(ctx, idempotencyInput.TenantID, idempotencyInput.Key, idempotencyInput.Payload)

	if err != nil {
		return nil, fmt.Errorf("failed to begin idempotency: %w", err)
	}

	if record.Status == "completed" {
		return s.repo.GetIntentByID(ctx, input.TenantID, record.ResponseRef)
	}



	intent := models.PaymentIntents{
		TenantID: input.TenantID,
		MerchantAccountID: input.MerchantAccountID,
		MerchantReference: input.MerchantReference,
		Amount: input.Amount,
		Currency: input.Currency,
		Status: "pending",
		ExpiresAt: input.ExpiresAt,
	}
	if err := s.repo.CreateIntent(ctx, &intent); err != nil{
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	if err := s.idempotencyService.Complete(ctx, input.TenantID, idempotencyInput.Key, intent.ID); err != nil{
		return nil, fmt.Errorf("failed to complete idempotency: %w", err)
	}

	return &intent, nil
}

func (s *Service) GetPaymentIntent(ctx context.Context, tenantID, id string ) (*models.PaymentIntents, error){
	if tenantID == "" || id == "" {
		return nil, fmt.Errorf("tenantID and id are required")
	}
	intent, err := s.repo.GetIntentByID(ctx, tenantID, id)
	if err != nil{
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return intent, nil
}

func (s *Service) ListPaymentIntens(ctx context.Context, tenantID string, limit, offset int) ([]*models.PaymentIntents, error){
	intents, err := s.repo.ListIntents(ctx, tenantID, limit, offset)
	if err != nil{
		return nil, fmt.Errorf("failed to list payment intents: %w", err)
	}
	return intents, nil
}

func (s *Service) SyncPaymentStatus(ctx context.Context, tenantID, paymentIntentID string) (*SyncPaymentResult, error){
	intent, err := s.repo.GetIntentByID(ctx, tenantID, paymentIntentID)

	if err != nil{
		return nil, fmt.Errorf("sync failed - intent not found: %w", err)
	}

	if intent.Status == "succeeded" || intent.Status == "expired" {
		return &SyncPaymentResult{PaymentIntent: intent}, nil
	}

	latestAttempt, err := s.repo.GetLatestAttempt(ctx, tenantID, paymentIntentID)
	if err != nil {
		if errors.Is(err, ErrNotFound){
			return &SyncPaymentResult{PaymentIntent: intent}, nil
		}
		return nil, fmt.Errorf("sync failed - get latest attempt: %w", err)
	}

	if latestAttempt.ProviderPaymentID == "" {
		return &SyncPaymentResult{PaymentIntent: intent, LatestAttempt: latestAttempt}, nil
	}

	providerClient, err := s.registry.GetProvider(ctx, latestAttempt.ProviderConnectionID)
	if err != nil {
		return nil, fmt.Errorf("sync failed - get provider: %w", err)
	}

	providerResp, err := providerClient.GetPaymentStatus(ctx, latestAttempt.ProviderPaymentID)

	if providerResp.NormalizedStatus != intent.Status{
		if err := s.repo.UpdateIntentStatus(ctx, tenantID, paymentIntentID, providerResp.NormalizedStatus);
		err != nil{
			return nil, fmt.Errorf("sync failed - update status: %w", err)
		}
		intent.Status = providerResp.NormalizedStatus
	}
	return &SyncPaymentResult{
		PaymentIntent: intent,
		LatestAttempt: latestAttempt,
		RawProviderStatus: providerResp.RawStatus,
	},nil
}
