package user

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/create_user"
)

var (
	passwordRegex = regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type CreateUserHandler interface {
	Handle(c *create_user.Command) error
}

type Controller struct {
	ge                *gin.Engine
	createUserHandler CreateUserHandler
}

func NewController(createUser CreateUserHandler) *Controller {
	return &Controller{
		createUserHandler: createUser,
	}
}

func (c *Controller) Register(ge *gin.Engine) {
	ge.POST("/v1/user", c.createUser)
	ge.POST("/v1/session", c.createSession) // TODO: authenticate with email and password
}

type Header struct {
	Value string `header:"Content-Type" example:"application/json" binding:"required"`
}

type CreateUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r CreateUserReq) Validate() error {
	if ok := emailRegex.MatchString(r.Email); !ok {
		return errors.New("invalid email")
	}

	if ok := passwordRegex.MatchString(r.Password); !ok {
		return errors.New("password must be at least 8 characters long and contain at least one letter and one number")
	}

	return nil
}

// createUser godoc
// @Summary Create a new user
// @Description Creates a new user account with an email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param Content-Type header string true "Content-Type" example(application/json)
// @Param data body CreateUserReq true "User data"
// @Success 201 {object} map[string]string "{"token": "token"}"
// @Failure 400 {object} map[string]string "{"error": "bad request"}"
// @Router /v1/user [post]
func (c *Controller) createUser(ctx *gin.Context) {
	var h Header
	err := ctx.ShouldBindHeader(&h)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var req CreateUserReq
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = req.Validate()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = c.createUserHandler.Handle(&create_user.Command{
		Email:    req.Email,
		Password: req.Password,
	})
}

type SessionHeader struct {
	Header
	Value string `header:"Authorization" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" binding:"required"`
}

// createSession godoc
// @Summary Create a new session
// @Description Creates a new session for a user by validating the email and password
// @Tags Sessions
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token" example("Bearer <your_jwt_token>")
// @Success 201 {object} map[string]string "{"token": "token"}"
// @Failure 400 {object} map[string]string "{"error": "invalid credentials"}"
// @Router /v1/session [post]
func (c *Controller) createSession(ctx *gin.Context) {

}
