package jwt

import (
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	jwtgo "github.com/lestrrat-go/jwx/jwt"
)

const ClientIDClaim string = "client_id"

type CreateClientToken struct {
	secret     string
	expiration time.Duration
}

func NewCreateClientToken(secret string) *CreateClientToken {
	return &CreateClientToken{secret: secret, expiration: time.Minute * 5}
}

func (c *CreateClientToken) Execute(publicID int64) (*string, error) {
	secretKey := []byte(c.secret)
	subStr := strconv.Itoa(int(publicID))

	token := jwtgo.New()
	_ = token.Set(jwtgo.IssuedAtKey, time.Now())
	_ = token.Set(jwtgo.ExpirationKey, time.Now().Add(c.expiration))
	_ = token.Set(jwtgo.SubjectKey, subStr)
	_ = token.Set(ClientIDClaim, subStr)
	_ = token.Set(jwtgo.NotBeforeKey, time.Now())

	signed, err := jwtgo.Sign(token, jwa.HS256, secretKey)
	if err != nil {
		return nil, err
	}
	jwtToken := string(signed)

	return &jwtToken, nil
}
