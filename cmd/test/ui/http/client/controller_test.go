package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/mail"
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
	passwdHahser *service.HashPassword
}

type clientResult struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	PublicID string `json:"public_id"`
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
	s.passwdHahser = hashPasswd

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

	const (
		clientEmail = "myman@myman.cz"
		clientPass  = "MyMan12345"
	)

	emailParsed, err := mail.ParseAddress(clientEmail)
	if err != nil {
		s.T().Fatal(err)
	}

	reqBody := map[string]string{
		"email":    clientEmail,
		"password": clientPass,
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

	row, cancel := s.pgConn.QueryRow(
		context.Background(),
		"TestGetClient",
		"SELECT id, email, public_id FROM client WHERE email = @email",
		pgx.NamedArgs{"email": emailParsed.String()},
	)
	defer cancel()

	clientRow := clientResult{}
	err = (*row).Scan(&clientRow.ID, &clientRow.Email, &clientRow.PublicID)
	if err != nil {
		s.T().Fatal(err)
	}
	s.Equal(emailParsed.String(), clientRow.Email)
	s.Equal(fmt.Sprintf("%d", decToken), clientRow.PublicID)

	s.T().Cleanup(func() {
		r, cancel, err := s.pgConn.Query(
			context.Background(),
			"TestDeleteClient",
			"DELETE FROM client WHERE id = @clientID",
			pgx.NamedArgs{"clientID": clientRow.ID},
		)
		if err != nil {
			s.T().Fatal(err)
		}
		defer cancel()

		if err := (*r).Err(); err != nil {
			s.T().Fatal(err)
		}
	})
}

func (s *ClientControllerTestSuite) Test_CreateClient_Fail_InvalidEmail() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	const (
		clientEmail = "myman@myman"
		clientPass  = "MyMan12345"
	)
	reqBody := map[string]string{
		"email":    clientEmail,
		"password": clientPass,
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

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ClientControllerTestSuite) Test_CreateClient_Fail_InvalidPassword() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	const (
		clientEmail = "myman@myman"
		clientPass  = "MyMan"
	)
	reqBody := map[string]string{
		"email":    clientEmail,
		"password": clientPass,
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

	s.Equal(http.StatusBadRequest, w.Code)
	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		s.T().Fatal(err)
	}
	s.Equal("password must be at least 8 characters long and contain at least one letter and one number", response["error"])
}

func (s *ClientControllerTestSuite) Test_CreateSession_Success() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	const (
		clientEmail = "myman@myman.cz"
		clientPass  = "MyMan12345"
	)

	emailParsed, err := mail.ParseAddress(clientEmail)
	if err != nil {
		s.T().Fatal(err)
	}

	row, cancel, err := s.pgConn.Query(
		context.Background(),
		"TestCreateClient",
		"INSERT INTO client (email, hashed_password) VALUES (@email, @password);",
		pgx.NamedArgs{"email": emailParsed.String(), "password": s.passwdHahser.Execute(clientPass)},
	)
	if err != nil {
		s.T().Fatal(err)
	}
	defer cancel()

	if err := (*row).Err(); err != nil {
		s.T().Fatal(err)
	}

	reqBody := map[string]string{
		"email":    clientEmail,
		"password": clientPass,
	}
	body, err := json.Marshal(&reqBody)
	if err != nil {
		s.T().Fatal(err)
	}
	r, _ := http.NewRequest(http.MethodPost, "/v1/session", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = r

	engine.Handle("POST", "/v1/session", s.clientCTRL.CreateSession)
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

	getClientRow, cancel := s.pgConn.QueryRow(
		context.Background(),
		"TestGetClient",
		"SELECT id, email, public_id FROM client WHERE email = @email",
		pgx.NamedArgs{"email": emailParsed.String()},
	)
	defer cancel()

	clientRow := clientResult{}
	err = (*getClientRow).Scan(&clientRow.ID, &clientRow.Email, &clientRow.PublicID)
	if err != nil {
		s.T().Fatal(err)
	}
	s.Equal(emailParsed.String(), clientRow.Email)
	s.Equal(fmt.Sprintf("%d", decToken), clientRow.PublicID)

	s.T().Cleanup(func() {
		r, cancel, err := s.pgConn.Query(
			context.Background(),
			"TestDeleteClient",
			"DELETE FROM client WHERE id = @clientID",
			pgx.NamedArgs{"clientID": clientRow.ID},
		)
		if err != nil {
			s.T().Fatal(err)
		}
		defer cancel()

		if err := (*r).Err(); err != nil {
			s.T().Fatal(err)
		}
	})
}

func (s *ClientControllerTestSuite) Test_CreateSession_Fail_WrongCredentials() {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	const (
		clientEmail = "myman@myman.cz"
		clientPass  = "MyMan12345"
	)

	reqBody := map[string]string{
		"email":    clientEmail,
		"password": clientPass,
	}
	body, err := json.Marshal(&reqBody)
	if err != nil {
		s.T().Fatal(err)
	}
	r, _ := http.NewRequest(http.MethodPost, "/v1/session", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = r

	engine.Handle("POST", "/v1/session", s.clientCTRL.CreateSession)
	engine.HandleContext(ctx)

	s.Equal(http.StatusUnauthorized, w.Code)
}

func TestClientControllerSuite(t *testing.T) {
	suite.Run(t, new(ClientControllerTestSuite))
}
