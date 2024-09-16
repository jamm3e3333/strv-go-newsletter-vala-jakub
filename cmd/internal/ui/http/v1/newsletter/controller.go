package newsletter

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/dto"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/query/list_newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/middleware"
)

type CreateNewsletterHandler interface {
	Handle(ctx context.Context, c *create_newsletter.Command) error
}

type ListNewsletterHandler interface {
	Execute(ctx context.Context, q list_newsletter.Query) ([]dto.Newsletter, error)
}

type Controller struct {
	createNewsletterHan CreateNewsletterHandler
	listNewsletterHan   ListNewsletterHandler
}

func NewController(createNewsletterHan CreateNewsletterHandler, listNewsletterHan ListNewsletterHandler) *Controller {
	return &Controller{
		createNewsletterHan: createNewsletterHan,
		listNewsletterHan:   listNewsletterHan,
	}
}

func (c *Controller) Register(ge *gin.Engine, authMiddleware gin.HandlerFunc) {
	ge.POST("/v1/newsletter", authMiddleware, c.createNewsletter)
	ge.GET("/v1/newsletter", authMiddleware, c.listNewsletter)
}

type Header struct {
	Value string `header:"Content-Type" example:"application/json" binding:"required"`
}

type CreateNewsletterReq struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// @Summary Create a new newsletter
// @Description Create a new newsletter with the specified name and optional description
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param Content-Type header string true "application/json"
// @Param CreateNewsletterReq body CreateNewsletterReq true "Newsletter details"
// @Success 201 {object} string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /v1/newsletter [post]
func (c *Controller) createNewsletter(ctx *gin.Context) {
	var h Header
	err := ctx.ShouldBindHeader(&h)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	clientID := ctx.GetInt64(middleware.ClientCTXKey)

	var req CreateNewsletterReq
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = c.createNewsletterHan.Handle(ctx, &create_newsletter.Command{
		ClientID:    clientID,
		Name:        req.Name,
		Description: req.Description,
	})

	// TODO: map error
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

type ListNewsletterQP struct {
	Page     int32 `form:"page" json:"page" binding:"required"`
	PageSize int32 `form:"page_size" json:"page_size" binding:"required"`
}

type ListNewsletterRes struct {
	Data       []Newsletter `json:"data"`
	TotalCount int64        `json:"total_count"`
}

type Newsletter struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ID          int64   `json:"id"`
	ClientID    int64   `json:"client_id"`
}

// @Summary List newsletters with pagination
// @Description Get a paginated list of newsletters
// @Tags Newsletter
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param page query int true "Page number" default(1)
// @Param pageSize query int true "Number of newsletters per page" default(10) minimum(5) maximum(100)
// @Success 200 {object} []CreateNewsletterReq "List of newsletters"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /v1/newsletter [get]
func (c *Controller) listNewsletter(ctx *gin.Context) {
	var q ListNewsletterQP
	err := ctx.ShouldBindQuery(&q)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	clientID := ctx.GetInt64(middleware.ClientCTXKey)
	newsletters, err := c.listNewsletterHan.Execute(ctx, list_newsletter.Query{
		ClientID: clientID,
		Page:     q.Page,
		PageSize: q.PageSize,
	})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(newsletters) <= 0 {
		ctx.JSON(http.StatusOK, ListNewsletterRes{
			Data:       []Newsletter{},
			TotalCount: 0,
		})
		return
	}

	newslettersRes := make([]Newsletter, len(newsletters))
	for i, v := range newsletters {
		newslettersRes[i] = Newsletter{
			Name:        v.Name,
			Description: v.Description,
			ID:          v.PublicID,
			ClientID:    v.ClientPublicID,
		}
	}

	ctx.JSON(http.StatusOK, ListNewsletterRes{
		Data:       newslettersRes,
		TotalCount: newsletters[0].TotalCount,
	})
}
