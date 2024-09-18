package mailjet

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type EmailClient struct {
	apiKey    string
	apiSecret string
	lg        logger.Logger
}

type Config struct {
	APIKey    string
	APISecret string
}

func NewEmailClient(lg logger.Logger, config Config) *EmailClient {
	return &EmailClient{
		apiKey:    config.APIKey,
		apiSecret: config.APISecret,
		lg:        lg,
	}
}

type Recipient struct {
	Email *mail.Address
	Name  string
}

type Sender struct {
	Email *mail.Address
	Name  string
}

type SendEmailParams struct {
	Recipient Recipient
	Sender    Sender
	Subject   string
	Text      string
	HTML      string
}

func (e *EmailClient) Send(ctx context.Context, p SendEmailParams) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mailjetClient := mailjet.NewMailjetClient(e.apiKey, e.apiSecret)
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: p.Sender.Email.String(),
				Name:  p.Sender.Name,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: p.Recipient.Email.String(),
					Name:  p.Recipient.Name,
				},
			},
			Subject:  p.Subject,
			TextPart: p.Text,
			HTMLPart: p.HTML,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		e.lg.ErrorWithMetadata("email sending failed", map[string]any{
			"recipient": p.Recipient.Email.String(),
			"subject":   p.Subject,
			"sender":    p.Sender.Email.String(),
			"text":      p.Text,
			"error":     err.Error(),
		})
		return err
	}

	e.lg.InfoWithMetadata("email sent", map[string]any{
		"recipient": p.Recipient.Email.String(),
		"subject":   p.Subject,
		"sender":    p.Sender.Email.String(),
		"text":      p.Text,
		"response":  fmt.Sprintf("%+v", res),
	})
	return nil
}
