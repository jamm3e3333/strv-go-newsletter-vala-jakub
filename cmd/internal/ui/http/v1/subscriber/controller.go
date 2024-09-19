package subscriber

import (
	"context"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_subscription"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/delete_subscription"
)

type CreateSubscriptionHandler interface {
	Handle(ctx context.Context, c create_subscription.Command) error
}

type DeleteSubscriptionHandler interface {
	Execute(ctx context.Context, c delete_subscription.Command) error
}

type Controller struct {
	createSubHan CreateSubscriptionHandler
	deleteSubHan DeleteSubscriptionHandler
}

func NewController(createSubHan CreateSubscriptionHandler, deleteSubHandler DeleteSubscriptionHandler) *Controller {
	return &Controller{createSubHan: createSubHan, deleteSubHan: deleteSubHandler}
}

func (c *Controller) Register(ge *gin.Engine) {
	ge.POST("/v1/newsletter/subscriber", c.createSubscription)
	ge.GET("/v1/newsletter/subscriber/unsubscribe", c.deleteSubscription)
}

type Header struct {
	Value string `header:"Content-Type" example:"application/json" binding:"required"`
}

type CreateSubscriptionReq struct {
	Email              string `json:"email" binding:"required"`
	NewsletterPublicID int64  `json:"newsletter_public_id" binding:"required"`
}

// createSubscription handles the subscription creation.
// @Summary      Create a new subscription
// @Description  Registers a new email subscription to a newsletter
// @Tags         Subscriber
// @Accept       json
// @Produce      json
// @Param        header body Header true "Content-Type header"
// @Param        body body CreateSubscriptionReq true "Subscription details"
// @Success      201  {string}  string  "Subscription created successfully"
// @Failure      400  {object}  gin.H   "Bad Request"
// @Failure      500  {object}  gin.H   "Internal Server Error"
// @Router       /v1/newsletter/subscriber [post]
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_SERVER_ERROR"})
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

type DeleteSubscriptionQP struct {
	NewsletterPublicID int64  `form:"newsletter_public_id" json:"newsletter_public_id" binding:"required"`
	Code               string `form:"code" json:"code" binding:"required"`
	Email              string `form:"email" json:"email" binding:"required"`
}

// deleteSubscription handles unsubscribing from a newsletter.
// @Summary      Unsubscribe from a newsletter
// @Description  Removes an email subscription using a verification code
// @Tags         Subscriber
// @Accept       json
// @Produce      json
// @Param        newsletter_public_id  query  int64  true  "Newsletter Public ID"  example(12345)
// @Param        code                  query  string true  "Verification Code"     example("ABC123")
// @Param        email                 query  string true  "Email address"         example("user@example.com")
// @Success      200  {string}  string  "Successfully unsubscribed"
// @Failure      400  {object}  gin.H   "Bad Request"
// @Failure      500  {object}  gin.H   "Internal Server Error"
// @Router       /v1/newsletter/subscriber/unsubscribe [get]
func (c *Controller) deleteSubscription(ctx *gin.Context) {
	var qp DeleteSubscriptionQP
	err := ctx.ShouldBindQuery(&qp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	emailParsed, err := mail.ParseAddress(qp.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.deleteSubHan.Execute(ctx, delete_subscription.Command{
		NewsletterPublicID: qp.NewsletterPublicID,
		Email:              emailParsed,
		VerificationCode:   qp.Code,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_SERVER_ERROR"})
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}
