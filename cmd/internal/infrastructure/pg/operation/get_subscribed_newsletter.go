package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type GetSubscribedNewsletterID struct {
	pgConn pgx.Connection
}

func NewGetSubscribedNewsletterIDOp(pgConn pgx.Connection) *GetSubscribedNewsletterID {
	return &GetSubscribedNewsletterID{
		pgConn: pgConn,
	}
}

type SubbedNewsletterRes struct {
	ID int64 `db:"newsletter_id"`
}

func (o *GetSubscribedNewsletterID) Execute(ctx context.Context, p dto.GetSubscribedNewsletter) (int64, error) {
	r, cancel := o.pgConn.QueryRow(ctx, "GetSubscribedNewsletterID", o.sql(), pgx.NamedArgs{
		"newsletterPublicID": p.NewsletterPublicID,
		"email":              p.Email,
		"code":               p.VerifCode,
	})
	defer cancel()

	res := SubbedNewsletterRes{}
	err := (*r).Scan(&res.ID)
	if err != nil {
		return -1, err
	}

	return res.ID, nil
}

func (o *GetSubscribedNewsletterID) sql() string {
	return `
SELECT
	n.id
FROM
	newsletter_subscribers ns
	INNER JOIN newsletter n ON n.id = ns.newsletter_id
WHERE
	n.public_id = @newsletterPublicID
	AND ns.email = @email
	AND unsubscribe_verification_code = @code;
`
}
