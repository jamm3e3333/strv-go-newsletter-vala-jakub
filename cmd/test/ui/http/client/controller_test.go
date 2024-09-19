package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/command/create_session"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/application/service"
	jwtinfra "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/jwt"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/pg/operation"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/http/v1/client"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/ui/jwt"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/test/helper"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
	"github.com/stretchr/testify/suite"
)

const (
	jwtSecret  = "secret"
	hashSecret = "secret"
)

type ClientControllerTestSuite struct {
	suite.Suite

	lg logger.Logger

	decryptToken *jwt.DecryptToken
	clientCTRL   *client.Controller
	pgConn       pgx.Connection
}

func (s *ClientControllerTestSuite) SetupSuite() {
	s.lg = helper.NewBlankLogger()

	pgCfg := helper.NewPostgresConfig()
	pgConn, err := pgx.NewConnectionPool(context.Background(), pgx.Config{
		ConnectionURL:     pgCfg.ConnectionURL(),
		LogLevel:          "info",
		MaxConnLifetime:   pgCfg.MaxConnLifetime(),
		MaxConnIdleTime:   pgCfg.MaxConnIdleTime(),
		QueryTimeout:      pgCfg.QueryTimeout(),
		DefaultMaxConns:   pgCfg.DefaultMaxConns(),
		DefaultMinConns:   pgCfg.DefaultMinConns(),
		HealthCheckPeriod: pgCfg.HealthCheckPeriod(),
	}, s.lg, helper.NewDummyMetrics())
	if err != nil {
		s.T().Fatal(err)
	}
	s.pgConn = pgConn

	createClient := operation.NewCreateClientOperation(s.pgConn)
	getClientData := operation.NewGetClientDataOperation(s.pgConn)

	hashPasswd := service.NewHashPassword(hashSecret)
	verifyPasswd := service.NewVerifyPassword(hashSecret)
	createToken := jwtinfra.NewCreateClientToken(jwtSecret)

	craeteClientHan := create_client.NewCreateClientHandler(hashPasswd, createToken, createClient)
	createSessionHan := create_session.NewCreateSessionHandler(verifyPasswd, getClientData, createToken)

	s.decryptToken = jwt.NewDecryptToken(jwtSecret)

	s.clientCTRL = client.NewController(craeteClientHan, createSessionHan)
}

func (s *ClientControllerTestSuite) TearDownSuite() {

}

func (s *ClientControllerTestSuite) Test_CreateClient_Success() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	reqBody := map[string]string{
		"email":    "myman@myman.cz",
		"password": "MyMan12345",
	}
	body, err := json.Marshal(&reqBody)
	if err != nil {
		s.T().Fatal(err)
	}
	r, _ := http.NewRequest(http.MethodPost, "/v1/client", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = r

	engine.Handle("POST", "/v1/client", s.clientCTRL.CreateClient)
	engine.HandleContext(ctx)

	s.Equal(http.StatusCreated, w.Code)

	var result map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		s.T().Fatal(err)
	}
	encToken, ok := result["token"].(string)
	if !ok {
		s.T().Fatal("Client token was not created")
	}

	s.Equal("string", reflect.TypeOf(encToken).Name(), "-------> assert token type ")

	decToken, err := s.decryptToken.Execute(encToken)
	if err != nil {
		s.T().Fatal(err)
	}

	fmt.Println(decToken)
}

func TestClientControllerSuite(t *testing.T) {
	suite.Run(t, new(ClientControllerTestSuite))
}
