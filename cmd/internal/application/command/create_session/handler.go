package create_session

import (
	"context"
	"errors"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type PasswordVerifier interface {
	Execute(password, hashedPassword string) bool
}

type GetClientDataOperation interface {
	GetForEmail(ctx context.Context, email string) (dto.ClientData, error)
}

type TokenEncrypter interface {
	Execute(publicID int64) (*string, error)
}

type CreateSessionHandler struct {
	verifyPassword PasswordVerifier
	getClientData  GetClientDataOperation
	createToken    TokenEncrypter
}

func NewCreateSessionHandler(verifyPassword PasswordVerifier, getClientData GetClientDataOperation, createToken TokenEncrypter) *CreateSessionHandler {
	return &CreateSessionHandler{
		verifyPassword: verifyPassword,
		getClientData:  getClientData,
		createToken:    createToken,
	}
}

func (h *CreateSessionHandler) Handle(ctx context.Context, c Command) (*string, error) {
	client, err := h.getClientData.GetForEmail(ctx, c.Email)
	if err != nil {
		return nil, err
	}

	if !h.verifyPassword.Execute(c.Password, client.HashedPassword) {
		return nil, errors.New("invalid password")
	}

	token, err := h.createToken.Execute(client.PublicID)
	if err != nil {
		return nil, err
	}

	return token, nil
}
