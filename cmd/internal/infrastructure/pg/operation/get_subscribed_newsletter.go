package operation

import (
	"context"

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

type subbedNewsletterRes struct {
	ID int64 `db:"id"`
}

func (o *GetSubscribedNewsletterID) GetForUserCode(ctx context.Context, email, code string) (int64, error) {
	r, cancel, err := o.pgConn.Query(ctx, "GetSubscribedNewsletterID", o.sql(), pgx.NamedArgs{
		"code":  code,
		"email": email,
	})
	if err != nil {
		return -1, err
	}
	defer cancel()

	if err := (*r).Err(); err != nil {
		return -1, err
	}

	var res subbedNewsletterRes
	err = (*r).Scan(&res.ID)
	if err != nil {
		return -1, err
	}

	return res.ID, nil
}

func (o *GetSubscribedNewsletterID) sql() string {
	return `
SELECT
	*
FROM
	newsletter_subscribers ns
WHERE
	ns.email = @email
	AND ns.unsubscribe_verification_code = @code;
`
}
