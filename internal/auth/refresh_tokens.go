package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const RefreshTokenDuration = 60*24*time.Hour

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte,32)

	if _,err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}