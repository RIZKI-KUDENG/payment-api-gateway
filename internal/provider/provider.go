package provider

import (
	"context"
	"net/http"
	"time"
)

type CreatePaymentRequest struct{
	PaymentIntentID string
	MerchantReference string
	Amount int64
	Currency string
	CustomerEmail string
	CustomerName string
	Description string
	ExpiresAt *time.Time
}

type CreatePaymentResponse struct{
	ProviderPaymentID string
	ProviderStatus string
	CheckoutURL string
	QRString string
	VirtualAccount string
}

type PaymentStatusResponse struct{
	ProviderPaymentID string
	RawStatus string
	NormalizedStatus string
	PaidAt *time.Time
}


type PaymentProvider interface{
	ProviderName() string
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error)
	GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatusResponse, error)
	VerifyWebhookSignature(req *http.Request, secret string) (bool, error)
}
