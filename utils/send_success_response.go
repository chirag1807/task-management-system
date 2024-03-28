package utils

import (
	"encoding/json"
	"net/http"
)

// SendSuccessResponse is a generalized function that return success response to client for any request
func SendSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
