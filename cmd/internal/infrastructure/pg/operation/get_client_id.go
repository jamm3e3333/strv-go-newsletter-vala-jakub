package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type GetClientID struct {
	pgConn pgx.Connection
}

func NewGetClientIDOperation(pgConn pgx.Connection) *GetClientID {
	return &GetClientID{
		pgConn: pgConn,
	}
}

type Result struct {
	ID int64 `json:"id"`
}

func (o *GetClientID) Execute(ctx context.Context, publicID int64) (int64, error) {

	r, cancel := o.pgConn.QueryRow(ctx, "GetClientID", o.sql(), pgx.NamedArgs{
		"publicID": publicID,
	})
	defer cancel()

	var res Result
	err := (*r).Scan(&res.ID)
	if err != nil {
		return -1, err
	}

	return res.ID, nil
}

func (o *GetClientID) sql() string {
	return `
SELECT
	id
FROM
	client
WHERE
	public_id = @publicID;
`
}
