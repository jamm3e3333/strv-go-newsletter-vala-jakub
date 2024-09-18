package dto

import "net/mail"

type SendSubConfirmation struct {
	RecipientEmailAddr *mail.Address
	Subject            string
	Text               string
	HTML               string
}
