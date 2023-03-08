package shorturlservice

import (
	"crypto/aes"
	"log"
	"math/rand"
)

var key = []byte("Increment #9 key")

// GenerateToken генератор случайной последовательности байт
func GenerateToken() ([]byte, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// CryptoToken шифрование данных
func CryptoToken(token []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	cryptoString := make([]byte, len(token))
	aesBlock.Encrypt(cryptoString, token)
	return cryptoString, nil
}

// DeCryptoToken расшифрование данных
func DeCryptoToken(token []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	deCryptoString := make([]byte, len(token))
	aesBlock.Decrypt(deCryptoString, token)
	return deCryptoString, nil
}
