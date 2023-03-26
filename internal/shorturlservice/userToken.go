package shorturlservice

import (
	"crypto/aes"
	"math/rand"
)

var key = []byte("Increment #9 key")

// GeneratorUser интерфейс генератор короткого URL
type GeneratorUser interface {
	GenerateToken() ([]byte, error)
}

// TestGeneratorUser структура генератора коротких URL для тестирования
type TestGeneratorUser struct {
	Result []byte
}

// GenerateToken замена генератора без случайных последовательностей
func (g TestGeneratorUser) GenerateToken() ([]byte, error) {
	return g.Result, nil
}

// RandomGeneratorUser структура генератора короткий  URL
type RandomGeneratorUser struct {
}

// GenerateToken генератор случайной последовательности байт
func (RandomGeneratorUser) GenerateToken() ([]byte, error) {
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
		return nil, err
	}
	deCryptoString := make([]byte, len(token))
	aesBlock.Decrypt(deCryptoString, token)
	return deCryptoString, nil
}
