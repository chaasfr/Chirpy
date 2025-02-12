package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/chaasfr/chirpy/internal/auth"
	"github.com/google/uuid"
)

type PolkaWebhookInput struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

type Data struct {
	UserID string `json:"user_id"`
}

const UserUpgradedEvent = "user.upgraded"

func (cfg *apiConfig) HandlerPolkaWebhook(rw http.ResponseWriter, req *http.Request) {
	receivedApiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("cannot extract api key from header %s",err)
		ReturnJsonError(rw, 401, "cannot get api key")
		return
	}

	if receivedApiKey != cfg.PolkaKey {
		ReturnJsonError(rw, 401, "wrong api key")
		return
	}
	
	var input PolkaWebhookInput
	GetInputStructFromJson(&input, rw, req)

	if input.Event != UserUpgradedEvent {
		rw.WriteHeader(204)
		return
	}

	userid, err := uuid.Parse(input.Data.UserID)
	if err != nil {
		log.Printf("error parsing uuid from polka %s", err)
		ReturnJsonError(rw, 404, "invalid data.user_id")
		return
	}

	err = cfg.DbQueries.UpgradeUser(req.Context(), userid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			ReturnJsonError(rw, 404, "user not found")
			return
		} else {
			log.Printf("error retrieving user %s", err)
			ReturnJsonGenericInternalError(rw)
			return
		}
	}

	rw.WriteHeader(204)
}
