package client

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_session"
)

var (
	passwordRegexPattern = `^[0-9a-zA-Z$&+,:;=?@#|'<>.\-^*()%!]{6,20}$`
	emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type CreateClientHandler interface {
	Handle(ctx context.Context, c *create_client.Command) (*string, error)
}

type CreateSessionHandler interface {
	Handle(ctx context.Context, c create_session.Command) (*string, error)
}

type Controller struct {
	createClientHandler  CreateClientHandler
	createSessionHandler CreateSessionHandler
}

func NewController(createClient CreateClientHandler, createSession CreateSessionHandler) *Controller {
	return &Controller{
		createClientHandler:  createClient,
		createSessionHandler: createSession,
	}
}

func (c *Controller) Register(ge *gin.Engine) {
	ge.POST("/v1/client", c.CreateClient)
	ge.POST("/v1/session", c.CreateSession)
}

type Header struct {
	Value string `header:"Content-Type" example:"application/json" binding:"required"`
}

type CreateClientReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateClientRes struct {
	Token string `json:"token"`
}

func (r CreateClientReq) ValidatePassword() error {
	passwordRegex := regexp.MustCompile(passwordRegexPattern)
	if ok := passwordRegex.MatchString(r.Password); !ok {
		return errors.New("password must be at least 8 characters long and contain at least one letter and one number")
	}

	return nil
}

func (r CreateClientReq) ValidateEmail() error {
	emailRegex := regexp.MustCompile(emailRegexPattern)
	if ok := emailRegex.MatchString(r.Email); !ok {
		return errors.New("invalid email")
	}

	return nil
}

// CreateClient godoc
// @Summary Create a new client
// @Description Creates a new client account with an email and password
// @Tags Client
// @Accept json
// @Produce json
// @Param Content-Type header string true "Content-Type" example(application/json)
// @Param data body CreateClientReq true "Client data"
// @Success 201 {object} CreateClientRes "{"token": "token"}"
// @Failure 400 {object} map[string]string "{"error": "bad request"}"
// @Failure 422 {object} map[string]string "{"error": "bad request"}"
// @Failure 500 {object} map[string]string "{"error": "bad request"}"
// @Router /v1/client [post]
func (c *Controller) CreateClient(ctx *gin.Context) {
	var h Header
	err := ctx.ShouldBindHeader(&h)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req CreateClientReq
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = req.ValidatePassword()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = req.ValidateEmail()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	email, err := mail.ParseAddress(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: add error mapping
	token, err := c.createClientHandler.Handle(ctx, &create_client.Command{
		Email:    email.String(),
		Password: req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_SERVER_ERROR"})
		return
	}

	ctx.JSON(http.StatusCreated, CreateClientRes{Token: *token})
}

type CreateSessionReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateSessionRes struct {
	Token string `json:"token"`
}

// CreateSession godoc
// @Summary Create a new session
// @Description Creates a new session for a user by validating the email and password
// @Tags Sessions
// @Accept json
// @Produce json
// @Param Content-Type header string true "Content-Type" example(application/json)
// @Param data body CreateSessionReq true "Client data"
// @Success 201 {object} map[string]string "{"token": "token"}"
// @Failure 400 {object} map[string]string "{"error": "invalid credentials"}"
// @Router /v1/session [post]
func (c *Controller) CreateSession(ctx *gin.Context) {
	var h Header
	err := ctx.ShouldBindHeader(&h)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req CreateSessionReq
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

	token, err := c.createSessionHandler.Handle(ctx, create_session.Command{
		Email:    email.String(),
		Password: req.Password,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, CreateSessionRes{Token: *token})
}
