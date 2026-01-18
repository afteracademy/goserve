package utility

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString(n uint64) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	str := hex.EncodeToString(b)
	return str, nil
}
