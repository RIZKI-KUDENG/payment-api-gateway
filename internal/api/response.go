package api

import (
	"encoding/json"
	"net/http"
)


type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorDetail struct{
	Code string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return
	}
}

func WriteSuccess(w http.ResponseWriter, status int, data any){
	WriteJSON(w, status, SuccessResponse{Data: data})
}

func WriteError(w http.ResponseWriter, status int, code string, message string){
	WriteJSON(w, status, ErrorResponse{Error: ErrorDetail{Code: code, Message: message}})
}
