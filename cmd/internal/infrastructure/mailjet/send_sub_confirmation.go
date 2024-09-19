package mailjet

import (
	"context"
	"net/mail"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/mailjet"
)

type SendSubConfirmation struct {
	emailSender       mailjet.MailClientSender
	senderMailAddress *mail.Address
	isLocalEnv        bool
}

func NewSendSubConfirmation(emailSender mailjet.MailClientSender, senderMail *mail.Address, isLocalEnv bool) *SendSubConfirmation {
	return &SendSubConfirmation{
		emailSender:       emailSender,
		senderMailAddress: senderMail,
		isLocalEnv:        isLocalEnv,
	}
}

func (s *SendSubConfirmation) Execute(ctx context.Context, p dto.SendSubConfirmation) error {
	if s.isLocalEnv {
		return nil
	}
	mailParams := mailjet.SendEmailParams{
		Recipient: mailjet.Recipient{
			Email: p.RecipientEmailAddr,
			Name:  p.RecipientEmailAddr.Name,
		},
		Sender: mailjet.Sender{
			Email: s.senderMailAddress,
			Name:  s.senderMailAddress.Name,
		},
		Subject: p.Subject,
		Text:    p.Text,
		HTML:    p.HTML,
	}
	err := s.emailSender.Send(ctx, mailParams)
	if err != nil {
		return err
	}
	return nil
}
