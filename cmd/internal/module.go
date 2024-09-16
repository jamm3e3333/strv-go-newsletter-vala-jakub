package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_session"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/query/list_newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/service"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/jwt"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/pg/operation"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/middleware"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/v1/client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/v1/newsletter"
	jwt2 "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/jwt"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type ModuleParams struct {
	HashSecret string
	JWTSecret  string
	PGConn     pgx.Connection
}

func RegisterModule(ge *gin.Engine, p ModuleParams) {
	hashPaswd := service.NewHashPassword(p.HashSecret)
	verifyPasswd := service.NewVerifyPassword(p.HashSecret)

	createToken := jwt.NewCreateClientToken(p.JWTSecret)

	getClientData := operation.NewGetClientData(p.PGConn)
	createClient := operation.NewCreateClientOperation(p.PGConn)
	createNewsletter := operation.NewCreateNewsletterOperation(p.PGConn)
	listNewsletter := operation.NewListNewsletterOperation(p.PGConn)

	createClientHan := create_client.NewCreateClientHandler(hashPaswd, createToken, createClient)
	createSessHan := create_session.NewCreateSessionHandler(verifyPasswd, getClientData, createToken)

	createNewsletterHan := create_newsletter.NewCreateNewsletterHandler(createNewsletter)
	listNewsletterHan := list_newsletter.NewListNewsletterHandler(listNewsletter)

	clientCTRL := client.NewController(createClientHan, createSessHan)
	newsletterCTRL := newsletter.NewController(createNewsletterHan, listNewsletterHan)

	decryptToken := jwt2.NewDecryptToken(p.JWTSecret)
	getClientID := operation.NewGetClientIDOperation(p.PGConn)
	authMiddleware := middleware.AuthMiddleware(getClientID, decryptToken)

	clientCTRL.Register(ge)
	newsletterCTRL.Register(ge, authMiddleware)
}
