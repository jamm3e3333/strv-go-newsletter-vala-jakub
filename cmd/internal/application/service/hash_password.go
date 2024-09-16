package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type HashPassword struct {
	secret string
}

func NewHashPassword(secret string) *HashPassword {
	return &HashPassword{secret: secret}
}

func (s *HashPassword) Execute(password string) string {
	secret := []byte(s.secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(password))

	return hex.EncodeToString(h.Sum(nil))
}
