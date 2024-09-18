package create_subscription

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/assets/template"
)

const (
	UnsubVerificationCodeLen int32 = 16
)

type CreateSubscriptionOperation interface {
	Execute(ctx context.Context, p dto.CreateSubscription) error
}

type GetNewsletterOperation interface {
	Execute(ctx context.Context, publicID int64) (dto.GetNewsletter, error)
}

type CreateNewsletterSubscriptionOperation interface {
	Execute(ctx context.Context, p dto.CreateNewsletterSubscription) error
}

type SendSubConfirmationOperation interface {
	Execute(ctx context.Context, p dto.SendSubConfirmation) error
}

type UnsubVerificationCodeGen = func(int32) (string, error)

type CreateSubscription struct {
	createSubscription           CreateSubscriptionOperation
	getNewsletterID              GetNewsletterOperation
	createNewsletterSubscription CreateNewsletterSubscriptionOperation
	unsubVerificationCodeGen     UnsubVerificationCodeGen
	sendSubConfirmation          SendSubConfirmationOperation
}

func NewCreateSubscriptionHandler(
	createSubscription CreateSubscriptionOperation,
	getNewsletterID GetNewsletterOperation,
	createNewsletterSubscription CreateNewsletterSubscriptionOperation,
	unsubVerificationCodeGen UnsubVerificationCodeGen,
	sendSubConfirmation SendSubConfirmationOperation,
) *CreateSubscription {
	return &CreateSubscription{
		createSubscription:           createSubscription,
		getNewsletterID:              getNewsletterID,
		createNewsletterSubscription: createNewsletterSubscription,
		unsubVerificationCodeGen:     unsubVerificationCodeGen,
		sendSubConfirmation:          sendSubConfirmation,
	}
}

func (h *CreateSubscription) Handle(ctx context.Context, c Command) error {
	newsletter, err := h.getNewsletterID.Execute(ctx, c.NewsletterPublicID)
	if err != nil {
		return err
	}

	verifCode, err := h.unsubVerificationCodeGen(UnsubVerificationCodeLen)
	if err != nil {
		return err
	}

	err = h.createNewsletterSubscription.Execute(ctx, dto.CreateNewsletterSubscription{
		Email:        c.Email.String(),
		NewsletterID: newsletter.ID,
		VerifCode:    verifCode,
	})
	if err != nil {
		return err
	}

	err = h.sendSubConfirmation.Execute(ctx, dto.SendSubConfirmation{
		RecipientEmailAddr: c.Email,
		Subject:            template.GetConfirmSubSubject(c.NewsletterPublicID),
		Text:               template.GetConfirmSubTxt(newsletter.Name, verifCode),
		HTML:               template.GetConfirmSubHTML(newsletter.Name, verifCode),
	})
	return h.createSubscription.Execute(ctx, dto.CreateSubscription{
		Email:              c.Email.String(),
		NewsletterPublicID: c.NewsletterPublicID,
	})
}
