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
	path := fmt.Sprintf("subscriber/%s", encodedEmail)

	err := c.fbConn.Create(ctx, "CreateSubscription", path, map[int64]any{
		p.NewsletterPublicID: true,
	})
	if err != nil {
		fmt.Println("err1", err.Error())
		return err
	}

	return nil
}
