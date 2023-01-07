package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/hex"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var hexCryproNewToken string
var storageUsers = "storageUsers.log"

func (s Server) serviceAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
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
			decryptoCookies, err := shorturlservice.DecryptoToken(deHexCookies)
			if err != nil {
				log.Fatal(err)
			}
			hexCookies := hex.EncodeToString(decryptoCookies)

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
		hexCryproNewToken = hex.EncodeToString(CryptoCookie)

		cookie := new(http.Cookie)
		cookie.Name = "Authentication"
		cookie.Value = hexCryproNewToken

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
