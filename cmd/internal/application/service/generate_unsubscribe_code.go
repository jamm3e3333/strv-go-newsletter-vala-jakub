package service

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateUnsubscribeCode(n int32) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
