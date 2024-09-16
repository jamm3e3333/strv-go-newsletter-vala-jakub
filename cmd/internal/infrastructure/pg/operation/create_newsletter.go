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

func (o *CreateNewsletter) Execute(ctx context.Context, p dto.CreateNewsletter) error {
	r, cancel, err := o.pgConn.Query(ctx, "CreateNewsletter", o.sql(), pgx.NamedArgs{
		"clientID":    p.ClientID,
		"name":        p.Name,
		"description": p.Description,
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

func (o *CreateNewsletter) sql() string {
	return `
INSERT INTO newsletter (client_id, name, description)
		values(@clientID, @name, @description);
`
}
