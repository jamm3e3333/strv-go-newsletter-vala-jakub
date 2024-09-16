package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	jwtgo "github.com/lestrrat-go/jwx/jwt"
)

var (
	InvalidTokenErr = errors.New("invalid token")
)

const ClientIDClaim string = "client_id"

type DecryptToken struct {
	secret string
}

func NewDecryptToken(secret string) *DecryptToken {
	return &DecryptToken{secret: secret}
}

func (d *DecryptToken) Execute(token string) (int64, error) {
	s := []byte(d.secret)
	t, err := jwtgo.ParseString(
		token,
		jwtgo.WithValidate(true),
		jwtgo.WithVerify(jwa.HS256, s),
		jwtgo.WithAcceptableSkew(time.Second*5),
	)
	if err != nil {
		return -1, err
	}

	clientIDUntyped, ok := t.Get(ClientIDClaim)
	if !ok {
		return -1, InvalidTokenErr
	}

	clientIDStr, ok := clientIDUntyped.(string)
	if !ok {
		return -1, InvalidTokenErr
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		return -1, err
	}

	return int64(clientID), nil
}
