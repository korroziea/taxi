package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var charsetLen = big.NewInt(int64(len(charset)))

const idLength = 24

func GenID() (string, error) {
	result := make([]byte, idLength)

	for i := range idLength {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
