package pg

import (
	"context"

	healthcheck "github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/health"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

type HealthIndicator struct {
	ctx  context.Context
	conn pgx.Connection
	lg   logger.Logger
}

func NewHealthIndicator(ctx context.Context, conn pgx.Connection, lg logger.Logger) *HealthIndicator {
	return &HealthIndicator{
		ctx:  ctx,
		conn: conn,
		lg:   lg,
	}
}

func (i *HealthIndicator) ComponentName() string {
	return "pg-strv-newsletter"
}

func (i *HealthIndicator) Status() healthcheck.Status {
	qr, cancel := i.conn.QueryRow(i.ctx, "health-status", "SELECT 1 AS ok", pgx.NamedArgs{})
	defer cancel()
	var ok int64
	err := (*qr).Scan(&ok)
	if err != nil {
		i.lg.Error("strv-newsletter postgres connection is down and threw %s!", err)

		return healthcheck.StatusDown
	}

	if ok == 0 {
		i.lg.Error("strv-newsletter postgres connection is down!")

		return healthcheck.StatusDown
	}

	return healthcheck.StatusUp
}
