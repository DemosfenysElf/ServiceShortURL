package router

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"ServiceShortURL/internal/shorturlservice"

	"github.com/labstack/echo"
)

var hexCryptoNewToken string
var storageUsers = "storageUsers.log"

func (s serverShortener) serviceAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	fmt.Println("==>> serviceAuthentication")
	return func(c echo.Context) error {
		requestCookies := c.Request().Cookies()

		consumerUser, err := shorturlservice.NewConsumer(storageUsers)
		if err != nil {
			log.Fatal(err)
		}
		defer consumerUser.Close()

		if (len(requestCookies) > 0) && (requestCookies[0].Name == "Authentication") {
			deHexCookies, err := hex.DecodeString(requestCookies[0].Value)
			if err != nil {
				log.Fatal(err)
			}
			deCryptoCookies, err := shorturlservice.DeCryptoToken(deHexCookies)
			if err != nil {
				log.Fatal(err)
			}
			hexCookies := hex.EncodeToString(deCryptoCookies)

			for {
				readUser, err := consumerUser.ReadUser()
				if err != nil {
					break
				}

				if hexCookies == readUser.ValueUser {
					shorturlservice.SetStructCookies("Authentication", hexCookies)
					return next(c)
				}
			}
		}

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
		hexCryptoNewToken = hex.EncodeToString(CryptoCookie)

		cookie := new(http.Cookie)
		cookie.Name = "Authentication"
		cookie.Value = hexCryptoNewToken

		shorturlservice.SetStructCookies("Authentication", hex.EncodeToString(newToken))

		producerUser, err := shorturlservice.NewProducer(storageUsers)
		if err != nil {
			log.Fatal(err)
		}
		defer producerUser.Close()

		if err := producerUser.WriteUser(shorturlservice.GetStructCookies()); err != nil {
			log.Fatal(err)
		}

		c.SetCookie(cookie)
		return next(c)
	}
}
