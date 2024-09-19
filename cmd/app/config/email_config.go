package config

import (
	"fmt"
	"net/mail"

	"github.com/ilyakaznacheev/cleanenv"
)

type EmailAddressField mail.Address

func (f *EmailAddressField) SetValue(s string) error {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return fmt.Errorf("can't parse email address: %w", err)
	}

	*f = EmailAddressField(*addr)
	return nil
}

type EmailConfig struct {
	SenderEmailAddress    string `env:"CONFIG_EMAIL_SENDER_EMAIL_ADDRESS"`
	APIKey                string `env:"CONFIG_EMAIL_API_KEY"`
	APISecret             string `env:"CONFIG_EMAIL_API_SECRET"`
	UnsubURL              string `env:"CONFIG_UNSUBSCRIBE_URL"`
	SenderEmailAddrParsed *mail.Address
}

func CreateEmailConfig() (EmailConfig, error) {
	var cfg EmailConfig
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	parsedEmail, err := mail.ParseAddress(cfg.SenderEmailAddress)
	if err != nil {
		return cfg, fmt.Errorf("can't parse email address: %w", err)
	}
	cfg.SenderEmailAddrParsed = parsedEmail

	return cfg, nil
}
