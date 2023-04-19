package router

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"ServiceShortURL/internal/shorturlservice"
)

var testStorageURL = "../test/shortsURl.log"
var testStorageUsers = "../test/storageUsers.log"

// для реализации авторизации пользователей в обход middleware при тестировании
var GeneratorUsers shorturlservice.GeneratorUser

func (s ServerShortener) kykiDB() *http.Cookie {
	consumerUser, err := shorturlservice.NewConsumer(testStorageUsers)
	if err != nil {
		fmt.Println(err)
	}
	defer consumerUser.Close()

	newToken, err := s.GeneratorUsers.GenerateToken()
	if err != nil {
		fmt.Println(err)
	}

	CryptoCookie, err := shorturlservice.CryptoToken(newToken)
	if err != nil {
		fmt.Println(err)
	}
	hexCryptoNewToken := hex.EncodeToString(CryptoCookie)

	shorturlservice.SetStructCookies("Authentication", hexCryptoNewToken)

	producerUser, err := shorturlservice.NewProducer(testStorageUsers)
	if err != nil {
		fmt.Println(err)
	}
	defer producerUser.Close()

	if err = producerUser.WriteUser(shorturlservice.GetStructCookies()); err != nil {
		fmt.Println(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authentication"
	cookie.Value = hexCryptoNewToken
	return cookie
}

func kyki() *http.Cookie {
	GeneratorUsers = shorturlservice.RandomGeneratorUser{}
	consumerUser, err := shorturlservice.NewConsumer(testStorageUsers)
	if err != nil {
		fmt.Println(err)
	}
	defer consumerUser.Close()

	newToken, err := GeneratorUsers.GenerateToken()
	if err != nil {
		fmt.Println(err)
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
			newToken, err = GeneratorUsers.GenerateToken()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			break
		}
	}

	CryptoCookie, err := shorturlservice.CryptoToken(newToken)
	if err != nil {
		fmt.Println(err)
	}
	hexCryptoNewToken := hex.EncodeToString(CryptoCookie)

	shorturlservice.SetStructCookies("Authentication", hex.EncodeToString(newToken))

	producerUser, err := shorturlservice.NewProducer(testStorageUsers)
	if err != nil {
		fmt.Println(err)
	}
	defer producerUser.Close()

	if err = producerUser.WriteUser(shorturlservice.GetStructCookies()); err != nil {
		fmt.Println(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authentication"
	cookie.Value = hexCryptoNewToken
	return cookie
}
