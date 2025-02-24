package postgres

import (
	"context"

	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/cmd/internal/infrastructure/prometheus"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/logger"
	"github.com/jamm3e3333/strv-go-newsletter-vala-jakub/pkg/pgx"
)

func EstablishConnection(ctx context.Context, cfg pgx.Config, lg logger.Logger, mm *prometheus.PgMetrics) *pgx.ConnectionPool {
	conn, err := pgx.NewConnectionPool(ctx, cfg, lg, mm.Cm)

	if err != nil {
		lg.Fatal("database connection error")
	}

	conn.RegisterMetrics(pgx.RegisterMetricsOptions{
		Qm: mm.Qm,
		Tm: mm.Tm,
	})

	return conn
}
