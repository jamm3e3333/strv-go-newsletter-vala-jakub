package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type GetNewsletter struct {
	pgConn pgx.Connection
}

func NewGetNewsletterOp(pgConn pgx.Connection) *GetNewsletter {
	return &GetNewsletter{
		pgConn: pgConn,
	}
}

type resultForNewsletter struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (o *GetNewsletter) Execute(ctx context.Context, publicID int64) (dto.GetNewsletter, error) {
	r, cancel := o.pgConn.QueryRow(ctx, "GetNewsletter", o.sql(), pgx.NamedArgs{
		"publicID": publicID,
	})
	defer cancel()
	res := resultForNewsletter{}

	err := (*r).Scan(&res.ID, &res.Name)
	if err != nil {
		return dto.GetNewsletter{}, err
	}

	return dto.GetNewsletter{
		ID:   res.ID,
		Name: res.Name,
	}, nil
}

func (o *GetNewsletter) sql() string {
	return `
SELECT
	id,
	name
FROM
	newsletter
WHERE
	public_id = @publicID;
`
}
