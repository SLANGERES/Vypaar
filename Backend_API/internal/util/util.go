package util

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/google/uuid"
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

func ParameterMissing(w http.ResponseWriter, StatusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	return json.NewEncoder(w).Encode(map[string]string{
		"sucess":  "false",
		"message": "Id is missing",
	})
}

func ParseInt(id string) (int64, error) {
	newId, err := strconv.ParseInt(id, 10, 64)

	return newId, err
}

func GenerateUID() string {
	token := uuid.NewString()

	return token
}
