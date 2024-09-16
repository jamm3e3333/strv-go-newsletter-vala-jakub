package list_newsletter

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type ListNewsletterOperation interface {
	Execute(ctx context.Context, p dto.ListNewsletter) ([]dto.Newsletter, error)
}

type ListNewsletter struct {
	listNewsletter ListNewsletterOperation
}

func NewListNewsletterHandler(listNewsletter ListNewsletterOperation) *ListNewsletter {
	return &ListNewsletter{
		listNewsletter: listNewsletter,
	}
}

func (o *ListNewsletter) Execute(ctx context.Context, q Query) ([]dto.Newsletter, error) {
	newsletters, err := o.listNewsletter.Execute(ctx, dto.ListNewsletter{
		ClientID: q.ClientID,
		Page:     q.Page,
		PageSize: q.PageSize,
	})
	if err != nil {
		return nil, err
	}

	return newsletters, nil
}
