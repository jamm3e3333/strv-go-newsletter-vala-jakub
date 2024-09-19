package firebase

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"google.golang.org/api/option"
)

type Connection struct {
	*db.Client

	lg logger.Logger
}

type Config struct {
	DBUrl      string
	SAKeyEnc   string
	IsTestMode bool
}

// NewConnection TODO: add metrics
func NewConnection(ctx context.Context, cfg Config, lg logger.Logger) *Connection {
	var app *firebase.App

	if cfg.IsTestMode {
		err := os.Setenv("FIREBASE_DATABASE_EMULATOR_HOST", "firebase-emulator:9000/?ns=strv-newsletter-go-vala-jakub-vala-local")
		if err != nil {
			lg.Fatal(err)
		}

		fbCfg := &firebase.Config{
			DatabaseURL: "https://strv-newsletter-go-vala-jakub-vala-local.firebaseio.com",
		}
		app, err = firebase.NewApp(ctx, fbCfg)
		if err != nil {
			lg.Fatal(err)
		}
	} else {
		saKey, err := base64.StdEncoding.DecodeString(cfg.SAKeyEnc)
		if err != nil {
			lg.Fatal(err)
		}
		opt := option.WithCredentialsJSON(saKey)
		fbCfg := &firebase.Config{
			DatabaseURL: cfg.DBUrl,
		}

		app, err = firebase.NewApp(ctx, fbCfg, opt)
		if err != nil {
			lg.Fatal(err)
		}
	}

	client, err := app.Database(ctx)
	if err != nil {
		lg.Fatal(err)
	}

	if err != nil {
		lg.Fatal(err)
	}

	return &Connection{
		client,
		lg,
	}
}

func (c *Connection) Create(ctx context.Context, opName, path string, data any) error {
	if err := c.NewRef(path).Set(ctx, data); err != nil {
		c.lg.ErrorWithMetadata(
			fmt.Sprintf("firebase op `%s` error", opName),
			map[string]any{
				"error": err.Error(),
				"name":  opName,
				"path":  path,
			},
		)
		return err
	}

	c.lg.InfoWithMetadata(fmt.Sprintf("success operation `%s`", opName), map[string]any{"path": path, "name": opName})
	return nil
}

func (c *Connection) Delete(ctx context.Context, opName, path string) error {
	err := c.NewRef(path).Delete(ctx)
	if err != nil {
		c.lg.ErrorWithMetadata(
			fmt.Sprintf("firebase op `%s` error", opName),
			map[string]any{
				"error": err.Error(),
				"name":  opName,
				"path":  path,
			},
		)
		return err
	}

	c.lg.InfoWithMetadata(fmt.Sprintf("success operation `%s`", opName), map[string]any{"path": path, "name": opName})
	return nil
}

func (c *Connection) GetForData(ctx context.Context, path, opName string, data any) error {
	if err := c.NewRef(path).Get(ctx, data); err != nil {
		c.lg.ErrorWithMetadata(
			fmt.Sprintf("firebase op `%s` error", opName),
			map[string]any{
				"error": err.Error(),
				"name":  opName,
				"path":  path,
			})
		return err

	}

	c.lg.InfoWithMetadata(
		fmt.Sprintf("success operation `%s`", opName),
		map[string]any{"data": data, "path": path, "name": opName},
	)
	return nil
}

func (c *Connection) Update(ctx context.Context, path, opName string, data map[string]any) error {
	err := c.NewRef(path).Update(ctx, data)
	if err != nil {
		c.lg.ErrorWithMetadata(
			fmt.Sprintf("firebase op `%s` error", opName),
			map[string]any{
				"error": err.Error(),
				"name":  opName,
				"path":  path,
			})
		return err
	}

	c.lg.InfoWithMetadata(
		fmt.Sprintf("success operation `%s`", opName),
		map[string]any{"data": data, "path": path, "name": opName},
	)
	return nil
}
