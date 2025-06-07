package util

import (
	"encoding/json"

	"net/http"

	"github.com/slangeres/Vypaar/backend_API/internal/types"
)

const StatusError = "404"

func WriteResponse(w http.ResponseWriter, StatusCode int, Message any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	return json.NewEncoder(w).Encode(Message)
}

func ErrorResponse(err error) types.ErrorResponse {
	var response types.ErrorResponse

	response.Status = StatusError
	response.Error = err.Error()

	return response
}
