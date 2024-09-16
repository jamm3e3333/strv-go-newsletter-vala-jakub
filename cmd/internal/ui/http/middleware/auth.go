package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ClientCTXKey = "clientID"

type GetClientIDOperation interface {
	Execute(ctx context.Context, publicID int64) (int64, error)
}

type ClientTokenDecrypter interface {
	Execute(token string) (int64, error)
}

func AuthMiddleware(getClientID GetClientIDOperation, decryptToken ClientTokenDecrypter) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") || len(authHeader) <= 7 {
			sendUnauthorized(c)
			return
		}
		bearerToken := authHeader[7:]

		publicClientID, err := decryptToken.Execute(bearerToken)
		if err != nil {
			sendUnauthorized(c)
			return
		}

		clientID, err := getClientID.Execute(c, publicClientID)

		c.Set(ClientCTXKey, clientID)
		c.Next()
	}
}

func sendUnauthorized(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
	c.Abort()
}
