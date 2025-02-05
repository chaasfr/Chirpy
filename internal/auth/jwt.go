package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const errorNoBearerToken = "bearer token not found in Authorization header"
const errorNoAuth =  "no authorization header provided"

const JwtDefaultDuration = time.Duration(1 * time.Hour)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			Subject:   userID.String(),
			Audience:  []string{},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expiresIn)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	)
	jwt, err :=token.SignedString([]byte(tokenSecret)) //we use byte because https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil //we use bytes because it is signed that way in MakeJWT
		},
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	uuidString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	
	uuID, err := uuid.Parse(uuidString)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New(errorNoAuth)
	}

	tokenString := ""
	authHeaderSplit := strings.Split(authHeader, " ")
	for i, word := range authHeaderSplit {
		if strings.ToLower(word) == "bearer" && i < len(authHeaderSplit) - 1 {
			tokenString = authHeaderSplit[i+1]
			return tokenString, nil
		}
	}

	return "", errors.New(errorNoBearerToken)
}