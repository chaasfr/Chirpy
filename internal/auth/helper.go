package auth

import (
	"errors"
	"net/http"
	"strings"
)

const errorNoKeyInAuth = " token not found in Authorization header"
const errorNoAuth = "no authorization header provided"

func GetAuthStringValue(headers http.Header, key string) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New(errorNoAuth)
	}

	tokenString := ""
	authHeaderSplit := strings.Split(authHeader, " ")
	lowerKey := strings.ToLower(key)
	for i, word := range authHeaderSplit {
		if strings.ToLower(word) == lowerKey && i < len(authHeaderSplit)-1 {
			tokenString = authHeaderSplit[i+1]
			return tokenString, nil
		}
	}

	return "", errors.New(key + errorNoKeyInAuth)
}
