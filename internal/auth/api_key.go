package auth

import (
	"net/http"
)

const ApiKeyKey = "apikey"

func GetAPIKey(headers http.Header) (string, error) {
	return GetAuthStringValue(headers, ApiKeyKey)
}