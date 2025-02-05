package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "amazingtestpasswordnotsafe"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password in %s", err)
		return
	}

	err = CheckPasswordHash(password, hashed)
	if err != nil {
		t.Errorf("CheckPasswordHash failed %s", err)
		return
	}
}