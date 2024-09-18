package delete_subscription

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type GetSubbedNewsletterIDOp interface {
	Execute(ctx context.Context, p dto.GetSubscribedNewsletter) (int64, error)
}

type DeleteSubscription struct {
	getSubbedNewsletterID GetSubbedNewsletterIDOp
}

func NewDeleteSubscription(getSubbedNewsletterID GetSubbedNewsletterIDOp) *DeleteSubscription {
	return &DeleteSubscription{
		getSubbedNewsletterID: getSubbedNewsletterID,
	}
}

func (h *DeleteSubscription) Execute(ctx context.Context, c Command) error {

}
