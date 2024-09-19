package delete_subscription

import "net/mail"

type Command struct {
	NewsletterPublicID int64
	Email              *mail.Address
	VerificationCode   string
}
