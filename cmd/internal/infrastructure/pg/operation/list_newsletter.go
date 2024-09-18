package operation

import (
	"context"
	"time"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type ListNewsletter struct {
	pgConn pgx.Connection
}

func NewListNewsletterOperation(pgConn pgx.Connection) *ListNewsletter {
	return &ListNewsletter{
		pgConn: pgConn,
	}
}

type newsletterResult struct {
	PublicID       int64     `db:"public_id"`
	Name           string    `db:"name"`
	ClientPublicID int64     `db:"client_public_id"`
	Description    *string   `db:"description"`
	TotalCount     int64     `db:"total_count"`
	CreatedAt      time.Time `db:"created_at"`
}

func (o *ListNewsletter) Execute(ctx context.Context, p dto.ListNewsletter) ([]dto.Newsletter, error) {
	offset := (p.Page - 1) * p.PageSize
	limit := p.PageSize

	r, cancel, err := o.pgConn.Query(ctx, "ListNewsletter", o.sql(), pgx.NamedArgs{
		"clientID": p.ClientID,
		"limit":    limit,
		"offset":   offset,
	})
	if err != nil {
		return nil, err
	}
	defer cancel()

	if err := (*r).Err(); err != nil {
		return nil, err
	}

	var newsletters []dto.Newsletter
	for (*r).Next() {
		var n newsletterResult
		err := (*r).Scan(
			&n.PublicID,
			&n.Description,
			&n.Name,
			&n.CreatedAt,
			&n.ClientPublicID,
			&n.TotalCount,
		)
		if err != nil {
			return nil, err
		}
		newsletters = append(newsletters, dto.Newsletter{
			ClientPublicID: n.ClientPublicID,
			PublicID:       n.PublicID,
			TotalCount:     n.TotalCount,
			Name:           n.Name,
			Description:    n.Description,
			CreatedAt:      n.CreatedAt,
		})
	}

	return newsletters, nil
}

func (o *ListNewsletter) sql() string {
	return `
SELECT
    nl.public_id,
    nl.description,
    nl.name,
    nl.created_at,
    cc.public_id AS client_public_id,
    COUNT(*) OVER() AS total_count
FROM
    client cc
    JOIN newsletter nl ON nl.client_id = cc.id
WHERE
    cc.id = @clientID
ORDER BY
    nl.created_at
LIMIT @limit OFFSET @offset;
`
}
