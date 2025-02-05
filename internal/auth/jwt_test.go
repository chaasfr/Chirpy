package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userid := uuid.New()
	tokenSecret := "sheeshIamASecret"
	wrongTokenSecret := "ohnoIamIncorrect"
	testDuration := time.Duration(1*time.Second)
	jwt, err := MakeJWT(userid,tokenSecret, testDuration)
	if err != nil {
		t.Errorf("error creating jwt: %s", err)
		return
	}

	idReceived, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Errorf("error validating jwt: %s", err)
		return
	}

	if idReceived !=userid {
		t.Errorf("error validating ids %v and %v", idReceived, userid)
	}

	_, err = ValidateJWT(jwt, wrongTokenSecret)
	if !strings.Contains(err.Error(), "signature is invalid") {
		t.Errorf("error validating jwt: %s", err)
		return
	}

	time.Sleep(testDuration)

	_, err = ValidateJWT(jwt, tokenSecret)
	if !strings.Contains(err.Error(), "token is expired") {
		t.Errorf("error validating jwt: %s", err)
		return
	}
}