package internal

import (
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_session"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_subscription"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/delete_subscription"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/query/list_newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/service"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/firebase"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/jwt"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/mailjet"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/pg/operation"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/middleware"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/v1/client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/v1/newsletter"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/v1/subscriber"
	jwtui "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/jwt"
	firebasepkg "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/firebase"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	emailclient "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/mailjet"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type ModuleParams struct {
	HashSecret      string
	JWTSecret       string
	EmailSenderAddr *mail.Address
	EmailUnsubURL   string
	AppENV          string
	PGConn          pgx.Connection
	FirebaseConn    firebasepkg.Connector
	Logger          logger.Logger
	MailClient      emailclient.MailClientSender
}

func RegisterModule(ge *gin.Engine, p ModuleParams) {
	hashPaswd := service.NewHashPassword(p.HashSecret)
	verifyPasswd := service.NewVerifyPassword(p.HashSecret)

	createToken := jwt.NewCreateClientToken(p.JWTSecret)

	getClientData := operation.NewGetClientDataOperation(p.PGConn)
	createClient := operation.NewCreateClientOperation(p.PGConn)
	createNewsletter := operation.NewCreateNewsletterOperation(p.PGConn)
	listNewsletter := operation.NewListNewsletterOperation(p.PGConn)
	getNewsletterID := operation.NewGetNewsletterOp(p.PGConn)
	createNewsletterSub := operation.NewCreateNewsletterSubscription(p.PGConn)
	isNewsletterSubExist := operation.NewIsNewsletterSubExistOp(p.PGConn)
	getSubbedNewsletterID := operation.NewGetSubscribedNewsletterIDOp(p.PGConn)
	delSub := firebase.NewDeleteSubscriptionOp(p.FirebaseConn)
	delNewsletterSub := operation.NewDeleteNewsletterSub(p.PGConn)

	createSub := firebase.NewCreateSubscription(p.FirebaseConn)
	createPubNewsletter := firebase.NewCreatePubNewsletter(p.FirebaseConn)

	sendSubConfirm := mailjet.NewSendSubConfirmation(
		p.MailClient,
		p.EmailSenderAddr,
		strings.Contains(p.AppENV, "local"),
	)

	createClientHan := create_client.NewCreateClientHandler(hashPaswd, createToken, createClient)
	createSessHan := create_session.NewCreateSessionHandler(verifyPasswd, getClientData, createToken)

	createNewsletterHan := create_newsletter.NewCreateNewsletterHandler(createNewsletter, createPubNewsletter)
	listNewsletterHan := list_newsletter.NewListNewsletterHandler(listNewsletter)
	createSubHan := create_subscription.NewCreateSubscriptionHandler(
		createSub,
		getNewsletterID,
		createNewsletterSub,
		service.GenerateUnsubscribeCode,
		sendSubConfirm,
		isNewsletterSubExist,
		p.EmailUnsubURL,
	)
	deleteSubHandler := delete_subscription.NewDeleteSubscriptionHandler(getSubbedNewsletterID, delSub, delNewsletterSub)

	clientCTRL := client.NewController(createClientHan, createSessHan)
	newsletterCTRL := newsletter.NewController(createNewsletterHan, listNewsletterHan)
	createSubCTRL := subscriber.NewController(createSubHan, deleteSubHandler)

	decryptToken := jwtui.NewDecryptToken(p.JWTSecret)
	getClientID := operation.NewGetClientIDOperation(p.PGConn)
	authMiddleware := middleware.AuthMiddleware(getClientID, decryptToken)

	clientCTRL.Register(ge)
	newsletterCTRL.Register(ge, authMiddleware)
	createSubCTRL.Register(ge)
}
