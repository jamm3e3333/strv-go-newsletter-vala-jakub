package mailjet

import "context"

type MailClientSender interface {
	Send(ctx context.Context, p SendEmailParams) error
}
