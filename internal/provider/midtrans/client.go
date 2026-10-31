package midtrans

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-api-gateway/internal/provider"
	"time"
)



type Client struct {
	serverKey string
	clientKey string
	environment string
	snapURL string
	apiURL string
	httpClient *http.Client
}


func NewClient(serverKey, clientKey, environment string) *Client{
	snapURL := "https://app.sandbox.midtrans.com"
    apiURL := "https://api.sandbox.midtrans.com"

    if environment == "production" {
    	snapURL = "https://app.midtrans.com"
        apiURL = "https://api.midtrans.com"
     }

    return &Client{
    serverKey: serverKey,
    clientKey: clientKey,
    snapURL: snapURL,
    apiURL: apiURL,
    httpClient:  &http.Client{Timeout: 10 * time.Second},
    }
}

func (c *Client) ProviderName() string{
	return "midtrans"
}

type snapTransactionRequest struct{
	TransactionDetails snapDetails `json:"transaction_details"`
}
type snapDetails struct{
	OrderID string `json:"order_id"`
	GrossAmount int64 `json:"gross_amount"`
}

type snapTransactionResponse struct {
	Token string `json:"token"`
	RedirectURL string `json:"redirect_url"`
	ErrorMessage []string `json:"error_message,omitempty"`
}

func (c *Client) CreatePayment(ctx context.Context, req provider.CreatePaymentRequest) (*provider.CreatePaymentResponse, error){

	endpoint := fmt.Sprintf("%s/snap/v1/transactions", c.snapURL)

	payload := snapTransactionRequest{
		TransactionDetails: snapDetails{
			OrderID: req.MerchantReference,
			GrossAmount: req.Amount,
		},
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("midtrans marshal payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))

	if err != nil {
		return nil, fmt.Errorf("midtrans create request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(c.serverKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+authHeader)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans http call: %w", err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("midtrans read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil,
		fmt.Errorf("midtrans error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var snapResp snapTransactionResponse
	if err := json.Unmarshal(respBytes, &snapResp); err != nil {
		return nil, fmt.Errorf("midtrans unmarshal response: %w", err)
	}

	return &provider.CreatePaymentResponse{
		ProviderPaymentID:
	 	snapResp.Token,
		ProviderStatus: "pending",
		CheckoutURL: snapResp.RedirectURL,
	}, nil
}

func (c *Client) GetPaymentStatus (ctx context.Context, orderID string)(*provider.PaymentStatusResponse, error){
	endpoint := fmt.Sprintf("%s/v2/%s/status", c.apiURL, orderID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(c.serverKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+authHeader)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result struct {
		TransactionStatus string `json:"transaction_status"`
		FraudStatus string `json:"fraud_status"`
		TransactionID string `json:"transaction_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	normalized := NormalizeStatus(result.TransactionStatus, result.FraudStatus)

	return &provider.PaymentStatusResponse{
		ProviderPaymentID: result.TransactionID,
		RawStatus: result.TransactionStatus,
		NormalizedStatus: normalized,
	}, nil
}

func (c *Client) VerifyWebhookSignature(req *http.Request, secret string) (bool, error){
	return true, nil
}
