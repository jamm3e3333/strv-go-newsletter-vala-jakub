package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type CreateClient struct {
	pgConn pgx.Connection
}

type CreateClientResult struct {
	PublicID int64 `db:"public_id"`
	ClientID int64 `db:"id"`
}

func NewCreateClientOperation(pgConn pgx.Connection) *CreateClient {
	return &CreateClient{pgConn: pgConn}
}

func (o *CreateClient) Execute(ctx context.Context, p dto.CreateClient) (dto.SavedClient, error) {
	r, cancel := o.pgConn.QueryRow(ctx, "CreateClient", o.sql(), pgx.NamedArgs{"email": p.Email, "hashedPassword": p.HashedPassword})
	defer cancel()

	res := CreateClientResult{}
	err := (*r).Scan(
		&res.PublicID,
		&res.ClientID,
	)
	if err != nil {
		return dto.SavedClient{}, err
	}

	return dto.SavedClient{
		ID:       res.ClientID,
		PublicID: res.PublicID,
	}, nil
}

func (o *CreateClient) sql() string {
	return `
INSERT INTO client (email, hashed_password)
		values(@email, @hashedPassword)
	RETURNING
		public_id, id;
`
}
