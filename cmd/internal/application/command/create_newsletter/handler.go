package create_newsletter

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type CreateNewsletterOperation interface {
	Execute(ctx context.Context, p dto.CreateNewsletter) (int64, error)
}

type CreatePublicNewsletterOp interface {
	Execute(ctx context.Context, p dto.CreatePublicNewsletter) error
}

type CreateNewsletter struct {
	createNewsletter    CreateNewsletterOperation
	createPubNewsletter CreatePublicNewsletterOp
}

func NewCreateNewsletterHandler(createNewsletter CreateNewsletterOperation, createPubNewsletter CreatePublicNewsletterOp) *CreateNewsletter {
	return &CreateNewsletter{
		createNewsletter:    createNewsletter,
		createPubNewsletter: createPubNewsletter,
	}
}

func (h *CreateNewsletter) Handle(ctx context.Context, c *Command) error {
	newsletterPubID, err := h.createNewsletter.Execute(ctx, dto.CreateNewsletter{
		Name:        c.Name,
		ClientID:    c.ClientID,
		Description: c.Description,
	})
	if err != nil {
		return err
	}

	err = h.createPubNewsletter.Execute(ctx, dto.CreatePublicNewsletter{
		Name:               c.Name,
		NewsletterPublicID: newsletterPubID,
	})
	if err != nil {
		return err
	}

	return nil
}
