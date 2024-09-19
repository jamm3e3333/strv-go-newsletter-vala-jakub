package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type IsNewsletterSubExist struct {
	pgConn pgx.Connection
}

func NewIsNewsletterSubExistOp(pgConn pgx.Connection) *IsNewsletterSubExist {
	return &IsNewsletterSubExist{
		pgConn: pgConn,
	}
}

type NewsletterSubCount struct {
	Count int64 `db:"count"`
}

func (o *IsNewsletterSubExist) Execute(ctx context.Context, p dto.GetNewsletterSub) (bool, error) {
	r, cancel := o.pgConn.QueryRow(ctx, "IsNewsletterSubExist", o.sql(), pgx.NamedArgs{
		"newsletterID": p.NewsletterID,
		"email":        p.Email,
	})
	defer cancel()

	res := NewsletterSubCount{}
	err := (*r).Scan(&res.Count)
	if err != nil {
		return false, err
	}
	return res.Count > 0, nil
}

func (o *IsNewsletterSubExist) sql() string {
	return `
SELECT
	count(*)
FROM
	newsletter_subscribers ns
WHERE
	ns.email = @email
	AND ns.newsletter_id = @newsletterID;
`
}
