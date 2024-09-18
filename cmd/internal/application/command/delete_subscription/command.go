package delete_subscription

type Command struct {
	NewsletterPublicID int64
	Email              string
	VerificationCode   string
}
