package operation

import (
	"context"
	"fmt"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type CreateNewsletterSubscription struct {
	pgConn pgx.Connection
}

func NewCreateNewsletterSubscription(pgConn pgx.Connection) *CreateNewsletterSubscription {
	return &CreateNewsletterSubscription{
		pgConn: pgConn,
	}
}

func (o *CreateNewsletterSubscription) Execute(ctx context.Context, p dto.CreateNewsletterSubscription) error {
	r, cancel, err := o.pgConn.Query(ctx, "CreateNewsletterSubscription", o.sql(), pgx.NamedArgs{
		"email":        p.Email,
		"newsletterID": p.NewsletterID,
		"code":         p.VerifCode,
	})
	if err != nil {
		fmt.Println("err1", err.Error())
		return err
	}
	defer cancel()

	if err := (*r).Err(); err != nil {
		fmt.Println("err2", err.Error())
		return err
	}

	return nil
}

func (o *CreateNewsletterSubscription) sql() string {
	return `
INSERT INTO newsletter_subscribers (email, newsletter_id, unsubscribe_verification_code)
		values(@email, @newsletterID, @code);
`
}
