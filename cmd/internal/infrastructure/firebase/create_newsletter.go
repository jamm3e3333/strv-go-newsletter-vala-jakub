package firebase

import (
	"context"
	"fmt"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/firebase"
)

type CreatePubNewsletter struct {
	fbConn firebase.Connector
}

func NewCreatePubNewsletter(fbConn firebase.Connector) *CreatePubNewsletter {
	return &CreatePubNewsletter{
		fbConn: fbConn,
	}
}

func (c *CreatePubNewsletter) Execute(ctx context.Context, p dto.CreatePublicNewsletter) error {
	err := c.fbConn.Create(ctx, "CreateNewsletter", "newsletter", map[string]any{
		fmt.Sprintf("%d", p.NewsletterPublicID): true,
	})
	if err != nil {
		return err
	}

	return nil
}
