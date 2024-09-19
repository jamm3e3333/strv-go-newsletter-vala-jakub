package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type DeleteNewsletterSub struct {
	pgConn pgx.Connection
}

func NewDeleteNewsletterSub(pgConn pgx.Connection) *DeleteNewsletterSub {
	return &DeleteNewsletterSub{
		pgConn: pgConn,
	}
}

func (o *DeleteNewsletterSub) Execute(ctx context.Context, p dto.DeleteNewsletterSubscription) error {
	r, cancel, err := o.pgConn.Query(ctx, "DeleteNewsletterSubscription", o.sql(), pgx.NamedArgs{
		"email":         p.Email,
		"newsletter_id": p.NewsletterID,
	})
	if err != nil {
		return err
	}
	defer cancel()

	if err := (*r).Err(); err != nil {
		return err
	}

	return nil
}

func (o *DeleteNewsletterSub) sql() string {
	return `
DELETE FROM newsletter_subscribers WHERE email = @email AND newsletter_id = @newsletter_id;
`
}
