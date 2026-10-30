package auth

import (
	"context"
	"errors"
	"uuid"
)



type SimpleAPIKeyAuthenticator struct{}


func (a *SimpleAPIKeyAuthenticator) Authenticate(ctx context.Context, apiKey string)(uuid.UUID, error){


	tenantID, err := uuid.Parse(apiKey)
	if err != nil {
		return uuid.UUID{}, errors.New("Invalid api key format")
	}

	return tenantID, nil
}
