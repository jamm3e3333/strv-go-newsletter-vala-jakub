package create_newsletter

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type CreateNewsletterOperation interface {
	Execute(ctx context.Context, p dto.CreateNewsletter) error
}

type CreateNewsletter struct {
	createNewsletter CreateNewsletterOperation
}

func NewCreateNewsletterHandler(createNewsletter CreateNewsletterOperation) *CreateNewsletter {
	return &CreateNewsletter{
		createNewsletter: createNewsletter,
	}
}

func (h *CreateNewsletter) Handle(ctx context.Context, c *Command) error {
	err := h.createNewsletter.Execute(ctx, dto.CreateNewsletter{
		Name:        c.Name,
		ClientID:    c.ClientID,
		Description: c.Description,
	})
	if err != nil {
		return err
	}

	return nil
}
