package create_client

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type UserPasswordHasher interface {
	Execute(password string) string
}

type TokenEncrypter interface {
	Execute(publicID int64) (*string, error)
}

type CreateClientOperation interface {
	Execute(ctx context.Context, p dto.CreateClient) (dto.SavedClient, error)
}

type CreateClient struct {
	hashPassword UserPasswordHasher
	createToken  TokenEncrypter
	createClient CreateClientOperation
}

func NewCreateClientHandler(hashPassword UserPasswordHasher, tokenEncrypter TokenEncrypter, createClient CreateClientOperation) *CreateClient {
	return &CreateClient{
		hashPassword: hashPassword,
		createToken:  tokenEncrypter,
		createClient: createClient,
	}
}

func (h *CreateClient) Handle(ctx context.Context, c *Command) (*string, error) {
	passwd := h.hashPassword.Execute(c.Password)
	client, err := h.createClient.Execute(ctx, dto.CreateClient{Email: c.Email, HashedPassword: passwd})
	if err != nil {
		return nil, err
	}

	token, err := h.createToken.Execute(client.PublicID)
	if err != nil {
		return nil, err
	}

	return token, err
}
