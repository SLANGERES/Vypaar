package util

import (
	"encoding/json"
	"strconv"

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

func ParseInt(id string) (int64, error) {
	newId, err := strconv.ParseInt(id, 10, 64)

	return newId, err
}
