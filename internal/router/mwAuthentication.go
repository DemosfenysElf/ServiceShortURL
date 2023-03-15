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

// MWAuthentication проверяем наличие ранее выданной куки
// если кука соответствует, то далее будет использоваться эта же кука
// если кука не соответствует или отсутствует, то выдаётся новая
func (s serverShortener) MWAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	fmt.Println("==>> MWAuthentication")
	return func(c echo.Context) error {
		requestCookies := c.Request().Cookies()

		consumerUser, err := shorturlservice.NewConsumer(storageUsers)
		if err != nil {
			log.Fatal(err)
		}
		defer consumerUser.Close()

		if (len(requestCookies) > 0) && (requestCookies[0].Name == "Authentication") {
			deHexCookies, errDecode := hex.DecodeString(requestCookies[0].Value)
			if errDecode != nil {
				log.Fatal(errDecode)
			}
			deCryptoCookies, errDecrypt := shorturlservice.DeCryptoToken(deHexCookies)
			if errDecrypt != nil {
				log.Fatal(errDecrypt)
			}
			hexCookies := hex.EncodeToString(deCryptoCookies)

			for {
				readUser, errRead := consumerUser.ReadUser()
				if errRead != nil {
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
			readUser, errRead := consumerUser.ReadUser()
			if errRead != nil {
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
