package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"google.golang.org/api/option"
)

type Connection struct {
	*db.Client

	lg logger.Logger
}

type Config struct {
	DBUrl      string
	SA         string
	IsTestMode bool
}

// NewConnection TODO: add metrics
func NewConnection(ctx context.Context, cfg Config, lg logger.Logger) *Connection {
	var opt option.ClientOption
	if cfg.IsTestMode {
		opt = option.WithoutAuthentication()
	} else {
		opt = option.WithCredentialsJSON([]byte(cfg.SA))
	}

	config := &firebase.Config{DatabaseURL: cfg.DBUrl}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		lg.Fatal(err)
	}
	client, err := app.Database(ctx)
	if err != nil {
		lg.Fatal(err)
	}

	return &Connection{
		client,
		lg,
	}
}

func (c *Connection) Create(ctx context.Context, opName, path string, data []byte) error {
	if err := c.NewRef(path).Set(ctx, data); err != nil {
		c.lg.ErrorWithMetadata(
			fmt.Sprintf("firebase op `%s` error", opName),
			map[string]any{
				"error": err.Error(),
				"data":  string(data),
				"name":  opName,
				"path":  path,
			},
		)
		return err
	}

	c.lg.InfoWithMetadata(fmt.Sprintf("success operation `%s`", opName), map[string]any{"data": string(data), "path": path, "name": opName})
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
