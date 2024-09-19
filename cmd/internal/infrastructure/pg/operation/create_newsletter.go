package operation

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type CreateNewsletter struct {
	pgConn pgx.Connection
}

func NewCreateNewsletterOperation(pgConn pgx.Connection) *CreateNewsletter {
	return &CreateNewsletter{
		pgConn: pgConn,
	}
}

type CreateNewsletterResult struct {
	PublicID int64 `db:"public_id"`
}

func (o *CreateNewsletter) Execute(ctx context.Context, p dto.CreateNewsletter) (int64, error) {
	r, cancel := o.pgConn.QueryRow(ctx, "CreateNewsletter", o.sql(), pgx.NamedArgs{
		"clientID":    p.ClientID,
		"name":        p.Name,
		"description": p.Description,
	})
	defer cancel()

	res := &CreateNewsletterResult{}
	err := (*r).Scan(&res.PublicID)
	if err != nil {
		return -1, err
	}

	return res.PublicID, nil
}

func (o *CreateNewsletter) sql() string {
	return `
INSERT INTO newsletter (client_id, name, description)
		values(@clientID, @name, @description) RETURNING public_id;
`
}
