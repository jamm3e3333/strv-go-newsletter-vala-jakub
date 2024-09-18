package dto

type CreateNewsletterSubscription struct {
	Email        string
	NewsletterID int64
	VerifCode    string
}
