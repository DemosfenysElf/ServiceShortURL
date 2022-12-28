package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/hex"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var userString string
var storageUsers = "storageUsers.log"

func (s Server) serviceAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestCookies := c.Request().Cookies()
		if (len(requestCookies) > 0) && (requestCookies[0].Name == "Authentication") {

			deHexCookies, err := hex.DecodeString(requestCookies[0].Value)
			if err != nil {
				log.Fatal(err)
			}
			decryptoCookies, err := shorturlservice.DecryptoToken(deHexCookies)
			if err != nil {
				log.Fatal(err)
			}

			consumerUser, err := shorturlservice.NewConsumer(storageUsers)
			if err != nil {
				log.Fatal(err)
			}
			defer consumerUser.Close()

			for {
				readUser, err := consumerUser.ReadUser()
				if err != nil {
					break
				}
				deHexUser, err := hex.DecodeString(readUser.ValueUser)
				if err != nil {
					log.Fatal(err)
				}

				if string(deHexUser) == string(decryptoCookies) {
					shorturlservice.SetStructCoockies("Authentication", readUser.ValueUser)
					return next(c)
				}
			}
		}

		consumerUser, err := shorturlservice.NewConsumer(storageUsers)
		if err != nil {
			log.Fatal(err)
		}
		defer consumerUser.Close()

		newToken, err := shorturlservice.GenerateToken()
		if err != nil {
			log.Fatal(err)
		}

		for {
			CryptoCookie, err := shorturlservice.CryptoToken(newToken)
			if err != nil {
				log.Fatal(err)
			}
			userString = hex.EncodeToString(CryptoCookie)

			readUser, err := consumerUser.ReadUser()
			if err != nil {
				break
			}
			deHexUser, err := hex.DecodeString(readUser.ValueUser)
			if err != nil {
				log.Fatal(err)
			}
			if string(deHexUser) != string(newToken) {
				break
			}
			newToken, err = shorturlservice.GenerateToken()
			if err != nil {
				log.Fatal(err)
			}
		}
		cookie := new(http.Cookie)
		cookie.Name = "Authentication"
		cookie.Value = userString

		shorturlservice.SetStructCoockies("Authentication", hex.EncodeToString(newToken))

		producerUser, err := shorturlservice.NewProducer(storageUsers)
		if err != nil {
			log.Fatal(err)
		}
		defer producerUser.Close()

		if err := producerUser.WriteUser(shorturlservice.GetStructCoockies()); err != nil {
			log.Fatal(err)
		}

		c.SetCookie(cookie)
		return next(c)
	}
}
