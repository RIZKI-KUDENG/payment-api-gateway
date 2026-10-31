package midtrans

import (
	"fmt"
	"net/http"
)

func GetBaseURL(environment string) (string, error){
	switch environment {
	case "production":
		return "https://app.midtrans.com", nil
	case "sandbox":
		return "https://app.sandbox.midtrans.com", nil
	default:
		return "", fmt.Errorf("invalid environment: %s", environment)
	}
}

func CreateTrxMidtrans(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	BaseURL, err := GetBaseURL(environment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
