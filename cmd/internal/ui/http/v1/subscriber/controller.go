package subscriber

import (
	"context"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_subscription"
)

type CreateSubscriptionHandler interface {
	Handle(ctx context.Context, c create_subscription.Command) error
}

type Controller struct {
	createSubHan CreateSubscriptionHandler
}

func NewController(createSubHan CreateSubscriptionHandler) *Controller {
	return &Controller{createSubHan: createSubHan}
}

func (c *Controller) Register(ge *gin.Engine) {
	ge.POST("/v1/subscriber", c.createSubscription)
}

type Header struct {
	Value string `header:"Content-Type" example:"application/json" binding:"required"`
}

type CreateSubscriptionReq struct {
	Email              string `json:"email" binding:"required"`
	NewsletterPublicID int64  `json:"newsletter_public_id" binding:"required"`
}

func (c *Controller) createSubscription(ctx *gin.Context) {
	var h Header
	err := ctx.ShouldBindHeader(&h)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req CreateSubscriptionReq
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := mail.ParseAddress(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.createSubHan.Handle(ctx, create_subscription.Command{
		Email:              email,
		NewsletterPublicID: req.NewsletterPublicID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}
