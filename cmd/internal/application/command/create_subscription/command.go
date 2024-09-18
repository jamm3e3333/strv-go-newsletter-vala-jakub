package create_subscription

import "net/mail"

type Command struct {
	Email              *mail.Address
	NewsletterPublicID int64
}
