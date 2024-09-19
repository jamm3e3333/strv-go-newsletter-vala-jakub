package delete_subscription

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
)

type GetSubbedNewsletterIDOp interface {
	Execute(ctx context.Context, p dto.GetSubscribedNewsletter) (int64, error)
}
type DeleteSubscriptionOperation interface {
	Execute(ctx context.Context, p dto.DeleteSubscription) error
}

type DeleteNewsletterSubscription interface {
	Execute(ctx context.Context, p dto.DeleteNewsletterSubscription) error
}

type DeleteSubscription struct {
	getSubbedNewsletterID GetSubbedNewsletterIDOp
	deleteSubscription    DeleteSubscriptionOperation
	deleteNewsletterSub   DeleteNewsletterSubscription
}

func NewDeleteSubscriptionHandler(getSubbedNewsletterID GetSubbedNewsletterIDOp, delSub DeleteSubscriptionOperation, delNewsletterSub DeleteNewsletterSubscription) *DeleteSubscription {
	return &DeleteSubscription{
		getSubbedNewsletterID: getSubbedNewsletterID,
		deleteSubscription:    delSub,
		deleteNewsletterSub:   delNewsletterSub,
	}
}

func (h *DeleteSubscription) Execute(ctx context.Context, c Command) error {
	newsletterID, err := h.getSubbedNewsletterID.Execute(ctx, dto.GetSubscribedNewsletter{
		Email:              c.Email.String(),
		NewsletterPublicID: c.NewsletterPublicID,
		VerifCode:          c.VerificationCode,
	})
	if err != nil {
		return err
	}

	err = h.deleteNewsletterSub.Execute(ctx, dto.DeleteNewsletterSubscription{
		Email:        c.Email.String(),
		NewsletterID: newsletterID,
	})
	if err != nil {
		return err
	}

	err = h.deleteSubscription.Execute(ctx, dto.DeleteSubscription{
		Email:              c.Email.String(),
		NewsletterPublicID: c.NewsletterPublicID,
	})
	if err != nil {
		return err
	}

	return nil
}
