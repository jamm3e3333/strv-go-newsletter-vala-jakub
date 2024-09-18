package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/app/config"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/app/setup/postgres"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/app/setup/prometheus"
	_ "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/app/swagger"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal"
	pghealth "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/pg"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/firebase"
	healthcheck "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/health"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/mailjet"
	pkgGin "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/net/http/gin"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/net/http/ginprometheus"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/net/http/server"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/shutdown"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

// Swagger API setup:
// @title STRV Newsletter API
// @version 2.0
// @description API provides endpoints for newsletter
// @contact.name STRV
func main() {
	ctx := shutdown.SetupShutdownContext()

	var (
		appConfig, errAPPConfig           = config.CreateAPPConfig()
		loggerConfig, errLoggerConfig     = config.CreateLoggerConfig()
		jwtConfig, errJWTConfig           = config.CreateJWTConfig()
		pgConfig, errPGConfig             = config.CreatePostgresConfig()
		firebaseConfig, errFirebaseConfig = config.CreateFirebaseConfig()
		emailConfig, errEmailConfig       = config.CreateEmailConfig()
	)

	for _, err := range []error{
		errAPPConfig,
		errLoggerConfig,
		errPGConfig,
		errJWTConfig,
		errFirebaseConfig,
		errEmailConfig,
	} {
		if err != nil {
			panic(err)
		}
	}

	location, _ := time.LoadLocation(appConfig.Timezone)
	time.Local = location

	lg := logger.New(logger.ParseLevel(loggerConfig.Level), loggerConfig.DevMode)

	mm := prometheus.NewMetricsOnce(appConfig.AppName)()
	pc := postgres.EstablishConnection(ctx, pgx.Config{
		ConnectionURL:     pgConfig.ConnectionURL(),
		LogLevel:          pgConfig.LogLevel,
		MaxConnLifetime:   pgConfig.MaxConnLifetime,
		MaxConnIdleTime:   pgConfig.MaxConnIdleTIme,
		QueryTimeout:      pgConfig.QueryTimeout,
		DefaultMaxConns:   pgConfig.MaxConns,
		DefaultMinConns:   pgConfig.MinConns,
		HealthCheckPeriod: pgConfig.HealthCheckPeriod,
	}, lg, mm.Pm)

	fbConn := firebase.NewConnection(
		ctx,
		firebase.Config{
			DBUrl:      firebaseConfig.DBUrl,
			SA:         firebaseConfig.SAKey,
			IsTestMode: firebaseConfig.IsTestingMode,
		},
		lg,
	)

	// Http server
	lg.Info("Initializing http server...")

	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	ge := gin.New()
	pprof.Register(ge)

	// Register logger middleware
	ge.Use(
		cors.New(cors.Config{
			AllowOrigins:     appConfig.AllowedOrigins(),
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowCredentials: true,
		}),
		pkgGin.LoggerMiddleware(pkgGin.NewLoggerMiddlewareConfig(
			[]string{"/metrics", "/health/liveness", "/health/readiness", "/status", "/api/*any"},
		), lg),
	)

	// Register prometheus endpoint and request/response metrics
	ge.GET("/metrics", ginprometheus.Handler())
	ge.Use(ginprometheus.Measure(ginprometheus.Config{
		Subsystem: appConfig.AppName,
		Labels:    []ginprometheus.Label{},
	}))
	lg.Info("gin prometheus initialized")

	// TODO: add recovery middleware

	// Initialize Swagger
	gsc := ginSwagger.Config{
		URL:                      "doc.json",
		DocExpansion:             "list",
		InstanceName:             swag.Name,
		Title:                    "STRV Newsletter API",
		DefaultModelsExpandDepth: 2,
		DeepLinking:              true,
		PersistAuthorization:     false,
		Oauth2DefaultClientID:    "",
	}
	ge.GET("/api/*any", ginSwagger.CustomWrapHandler(&gsc, swaggerFiles.Handler))
	lg.Info("swagger initialized")

	// Initialize health check
	readinessHealthCheck := healthcheck.NewHealthCheck(appConfig.HealthCheckTimeout, lg)
	_ = healthcheck.NewHealthCheck(appConfig.HealthCheckTimeout, lg)
	lg.Info("health checks initialized")

	readinessHealthCheck.RegisterIndicator(pghealth.NewHealthIndicator(ctx, pc, lg))
	lg.Info("health check indicators registered")

	mailClient := mailjet.NewEmailClient(lg, mailjet.Config{
		APIKey:    emailConfig.APIKey,
		APISecret: emailConfig.APISecret,
	})
	internal.RegisterModule(ge, internal.ModuleParams{
		HashSecret:      appConfig.HashSecret,
		JWTSecret:       jwtConfig.Secret,
		EmailSenderAddr: emailConfig.SenderEmailAddrParsed,
		PGConn:          pc,
		FirebaseConn:    fbConn,
		Logger:          lg,
		MailClient:      mailClient,
	})

	for _, v := range ge.Routes() {
		lg.Info("[HTTP] Route: %s %s initialized.", v.Method, v.Path)
	}
	lg.Info("Internal module initialized.")
	lg.Info("[HTTP] Gin initialized.")

	srv := server.NewServer(ge, appConfig.ReadTimeout, appConfig.WriteTimeout, appConfig.Port, appConfig.ShutdownTimeout)
	lg.Info("[HTTP] Server initialized.")

	lg.Info("[HTTP] Start listening on port %d.", appConfig.Port)
	httpErrChan := srv.Run()

	select {
	case err := <-httpErrChan:
		lg.Error("http server error, %s", err)
		shutdown.SignalShutdown()
	case <-ctx.Done():
		if err := srv.Shutdown(ctx); err != nil {
			lg.Error("err shutting down http server, error: %v", err)
		}
		lg.Info("shutdown signaled")
	}
}
