package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userid := uuid.New()
	tokenSecret := "sheeshIamASecret"
	wrongTokenSecret := "ohnoIamIncorrect"
	testDuration := time.Duration(1 * time.Second)
	jwt, err := MakeJWT(userid, tokenSecret, testDuration)
	if err != nil {
		t.Errorf("error creating jwt: %s", err)
		return
	}

	idReceived, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Errorf("error validating jwt: %s", err)
		return
	}

	if idReceived != userid {
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

func TestGetBearerToken(t *testing.T) {
	goodToken := "good_token!"
	goodHeader := http.Header{}
	badHeaderNoAuth := http.Header{}
	badHeaderNoBearer := http.Header{}

	goodHeader.Add("Authorization", "random stuff and Bearer "+goodToken)
	badHeaderNoAuth.Add("Content-Type", "text/json")
	badHeaderNoBearer.Add("Authorization", "some random bs")

	token, err := GetBearerToken(goodHeader)
	if err != nil {
		t.Errorf("error getting bearer token: %s", err)
		return
	}

	if token != goodToken {
		t.Errorf("error unexpected token %s instead of %s", token, goodToken)
		return
	}

	_, err = GetBearerToken(badHeaderNoAuth)
	if err.Error() != errorNoAuth {
		t.Errorf("wrong error when no auth: %s", err)
		return
	}

	_, err = GetBearerToken(badHeaderNoBearer)
	if err.Error() != BearerKey+errorNoKeyInAuth {
		t.Errorf("wrong error when no auth: %s", err)
		return
	}
}
