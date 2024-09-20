package firebase

import (
	"context"
	"fmt"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/firebase"
)

type CreateSubscription struct {
	fbConn firebase.Connector
}

func NewCreateSubscription(fbConn firebase.Connector) *CreateSubscription {
	return &CreateSubscription{
		fbConn: fbConn,
	}
}

func (c *CreateSubscription) Execute(ctx context.Context, p dto.CreateSubscription) error {
	encodedEmail := encodeEmail(p.Email)
	path := fmt.Sprintf("subscriber/%s/newsletter", encodedEmail)
	err := c.fbConn.Create(ctx, "CreateSubscription", path, map[string]any{
		fmt.Sprintf("%d", p.NewsletterPublicID): true,
	})
	if err != nil {
		return err
	}

	return nil
}
