package test

import (
	"encoding/hex"
	"log"
	"net/http"

	"ServiceShortURL/internal/shorturlservice"
)

type shortURLApiShortenBatch struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

func kyki() *http.Cookie {

	consumerUser, err := shorturlservice.NewConsumer("storageUsers.log")
	if err != nil {
		log.Fatal(err)
	}
	defer consumerUser.Close()

	newToken, err := shorturlservice.GenerateToken()
	if err != nil {
		log.Fatal(err)
	}

	userMap := make(map[string]bool, 0)

	for {
		readUser, err := consumerUser.ReadUser()
		if err != nil {
			break
		}
		userMap[readUser.ValueUser] = true
	}

	for {
		hexNewToken := hex.EncodeToString(newToken)
		if _, ok := userMap[hexNewToken]; ok {
			newToken, err = shorturlservice.GenerateToken()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			break
		}
	}

	CryptoCookie, err := shorturlservice.CryptoToken(newToken)
	if err != nil {
		log.Fatal(err)
	}
	hexCryptoNewToken := hex.EncodeToString(CryptoCookie)

	shorturlservice.SetStructCookies("Authentication", hex.EncodeToString(newToken))

	producerUser, err := shorturlservice.NewProducer("storageUsers.log")
	if err != nil {
		log.Fatal(err)
	}
	defer producerUser.Close()

	if err := producerUser.WriteUser(shorturlservice.GetStructCookies()); err != nil {
		log.Fatal(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authentication"
	cookie.Value = hexCryptoNewToken
	return cookie
}
