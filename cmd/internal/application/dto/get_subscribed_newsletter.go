package dto

type GetSubscribedNewsletter struct {
	Email              string
	NewsletterPublicID int64
	VerifCode          string
}
