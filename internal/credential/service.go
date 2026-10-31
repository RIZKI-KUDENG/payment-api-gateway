package credential

import (
	"context"
	"encoding/json"
	"fmt"
)



type Service struct {
	repo Repository
	masterKey []byte
}


func NewService(repo Repository, masterKey []byte) *Service{
	return &Service{repo: repo, masterKey: masterKey}
}

func (s *Service) StoreMidtransCredential(
	ctx context.Context,
	tenantID,
	merchantAccountId string,
	serverKey,
	clientKey,
	merchantID string,
) (string, error){

	data, err := json.Marshal(map[string]string{
		"server_key": serverKey,
		"client_key": clientKey,
	})

	if err != nil {
		return "",fmt.Errorf("marshal credential: %w", err)
	}

	encryptedVal, err := Encrypt(data, s.masterKey)

	if err != nil {
		return "", fmt.Errorf("encrypt credential: %w", err)
	}

	credentialID, err := s.repo.SaveCredential(ctx, tenantID, encryptedVal, 1)

	if err != nil {
		return "", err
	}

	return s.repo.CreateProviderConnection(
		ctx,
		ProviderConnectionInput{
			TenantID: tenantID,
			MerchantAccountID: merchantAccountId,
			Provider: "midtrans",
			ProviderAccountID: merchantID,
			Name: "Midtrans Sandbox",
			CredentialsRef: credentialID,
			Environment: "sandbox",
		},
	)
}
