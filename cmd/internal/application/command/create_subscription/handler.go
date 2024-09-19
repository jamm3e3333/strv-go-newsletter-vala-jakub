package create_subscription

import (
	"context"
	"fmt"
	"net/url"

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

type IsNewsletterSubExistOperation interface {
	Execute(ctx context.Context, p dto.GetNewsletterSub) (bool, error)
}

type UnsubVerificationCodeGen = func(int32) (string, error)

type CreateSubscription struct {
	createSubscription           CreateSubscriptionOperation
	getNewsletterID              GetNewsletterOperation
	createNewsletterSubscription CreateNewsletterSubscriptionOperation
	unsubVerificationCodeGen     UnsubVerificationCodeGen
	sendSubConfirmation          SendSubConfirmationOperation
	isNewsletterSubExist         IsNewsletterSubExistOperation
	unsubURL                     string
}

func NewCreateSubscriptionHandler(
	createSubscription CreateSubscriptionOperation,
	getNewsletterID GetNewsletterOperation,
	createNewsletterSubscription CreateNewsletterSubscriptionOperation,
	unsubVerificationCodeGen UnsubVerificationCodeGen,
	sendSubConfirmation SendSubConfirmationOperation,
	isNewsletterSubExist IsNewsletterSubExistOperation,
	unsubURL string,
) *CreateSubscription {
	return &CreateSubscription{
		createSubscription:           createSubscription,
		getNewsletterID:              getNewsletterID,
		createNewsletterSubscription: createNewsletterSubscription,
		unsubVerificationCodeGen:     unsubVerificationCodeGen,
		sendSubConfirmation:          sendSubConfirmation,
		isNewsletterSubExist:         isNewsletterSubExist,
		unsubURL:                     unsubURL,
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

	isExist, err := h.isNewsletterSubExist.Execute(ctx, dto.GetNewsletterSub{
		NewsletterID: newsletter.ID,
		Email:        c.Email.String(),
	})
	if err != nil {
		return err
	}

	if !isExist {
		err = h.createNewsletterSubscription.Execute(ctx, dto.CreateNewsletterSubscription{
			Email:        c.Email.String(),
			NewsletterID: newsletter.ID,
			VerifCode:    verifCode,
		})
		if err != nil {
			return err
		}

		params := url.Values{}
		params.Add("email", c.Email.Address)
		params.Add("code", verifCode)
		params.Add("newsletter_public_id", fmt.Sprintf("%d", c.NewsletterPublicID))
		unsubLink := fmt.Sprintf("%s?%s", h.unsubURL, params.Encode())

		fmt.Println("unsub_link", unsubLink)
		err = h.sendSubConfirmation.Execute(ctx, dto.SendSubConfirmation{
			RecipientEmailAddr: c.Email,
			Subject:            template.GetConfirmSubSubject(c.NewsletterPublicID),
			Text:               template.GetConfirmSubTxt(newsletter.Name, unsubLink),
			HTML:               template.GetConfirmSubHTML(newsletter.Name, unsubLink),
		})
		if err != nil {
			fmt.Println("error sending confirmation email:", err.Error())
		}

		err = h.createSubscription.Execute(ctx, dto.CreateSubscription{
			Email:              c.Email.String(),
			NewsletterPublicID: c.NewsletterPublicID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
