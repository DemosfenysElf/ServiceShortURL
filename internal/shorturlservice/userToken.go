package shorturlservice

import (
	"crypto/aes"
	"log"
	"math/rand"
)

var key = []byte("Increment #9 key")

func GenerateToken() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func CryptoToken(token []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	cryptStroka := make([]byte, len(token))
	aesblock.Encrypt(cryptStroka, token)
	return cryptStroka, nil
}

func DecryptoToken(token []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	deCryptStroka := make([]byte, len(token))
	aesblock.Decrypt(deCryptStroka, token)
	return deCryptStroka, nil
}
