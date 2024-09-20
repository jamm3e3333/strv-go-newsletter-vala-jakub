package firebase

import (
	"context"
	"fmt"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/firebase"
)

type DeleteSubscription struct {
	fbConn firebase.Connector
}

func NewDeleteSubscriptionOp(fbConn firebase.Connector) *DeleteSubscription {
	return &DeleteSubscription{
		fbConn: fbConn,
	}
}

func (c *DeleteSubscription) Execute(ctx context.Context, p dto.DeleteSubscription) error {
	encodedEmail := encodeEmail(p.Email)
	path := fmt.Sprintf("subscriber/%s/newsletter/%d", encodedEmail, p.NewsletterPublicID)
	err := c.fbConn.Delete(ctx, "DeleteSubscription", path)
	if err != nil {
		return err
	}

	return nil
}
