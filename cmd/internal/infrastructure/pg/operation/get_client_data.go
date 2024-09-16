package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type GetClientData struct {
	pgConn pgx.Connection
}

func NewGetClientData(pgConn pgx.Connection) *GetClientData {
	return &GetClientData{pgConn: pgConn}
}

type result struct {
	Email          string `db:"email"`
	HashedPassword string `db:"hashed_password"`
	PublicID       int64  `db:"public_id"`
}

func (o *GetClientData) GetForEmail(ctx context.Context, email string) (dto.ClientData, error) {
	r, cancel := o.pgConn.QueryRow(ctx, "GetClientData", o.sql(), pgx.NamedArgs{"email": email})
	defer cancel()

	res := result{}
	err := (*r).Scan(
		&res.Email,
		&res.HashedPassword,
		&res.PublicID,
	)
	if err != nil {
		return dto.ClientData{}, err
	}

	return dto.ClientData{
		Email:          res.Email,
		HashedPassword: res.HashedPassword,
		PublicID:       res.PublicID,
	}, nil
}

func (o *GetClientData) sql() string {
	return `
SELECT
	email,
	hashed_password,
	public_id
FROM
	client
WHERE
	email = @email;
`
}
