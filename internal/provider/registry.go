package provider

import "context"


type Registry interface {
	GetProvider(ctx context.Context, connectionID string) (PaymentProvider, error)
}
