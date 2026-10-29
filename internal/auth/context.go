package auth

import (
	"context"
	"errors"
	"uuid"
)


type contextKey string


const tenantIDKey contextKey = "tenant_id"


var ErrTenantIDNotFound = errors.New("tenant_id not found")


func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context{
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func GetTenantID(ctx context.Context) (uuid.UUID, error) {
	tenantID, ok := ctx.Value(tenantIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrTenantIDNotFound
	}
	return tenantID, nil
}
