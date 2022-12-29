package router

import (
	"ServiceShortURL/internal/shorturlservice"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

var hexCryproNewToken string
var storageUsers = "storageUsers.log"

func (s Server) serviceAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestCookies := c.Request().Cookies()
		fmt.Println(">>>>>REQUEST<<<<<<", requestCookies)
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

				if hexCookies == readUser.ValueUser {
					shorturlservice.SetStructCoockies("Authentication", hexCookies)
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
			readUser, err := consumerUser.ReadUser()
			if err != nil {
				break
			}
			hexNewToken := hex.EncodeToString(newToken)

			if hexNewToken != readUser.ValueUser {
				break
			}
			newToken, err = shorturlservice.GenerateToken()
			if err != nil {
				log.Fatal(err)
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

		shorturlservice.SetStructCoockies("Authentication", hex.EncodeToString(newToken))

		producerUser, err := shorturlservice.NewProducer(storageUsers)
		if err != nil {
			log.Fatal(err)
		}
		defer producerUser.Close()

		if err := producerUser.WriteUser(shorturlservice.GetStructCoockies()); err != nil {
			log.Fatal(err)
		}
		fmt.Println(">>>>>REQUEST-03<<<<<<", requestCookies)
		c.SetCookie(cookie)
		return next(c)
	}
}
