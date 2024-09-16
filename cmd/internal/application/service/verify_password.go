package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type VerifyPassword struct {
	secret string
}

func NewVerifyPassword(secret string) *VerifyPassword {
	return &VerifyPassword{secret: secret}
}

func (s *VerifyPassword) Execute(password, hashedPassword string) bool {
	secret := []byte(s.secret)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(password))

	hash := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(hashedPassword), []byte(hash))
}
