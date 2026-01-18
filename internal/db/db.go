package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"initialservice/internal/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.PgConn) (DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s ",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName,
	)

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return DB{}, fmt.Errorf("pgxpool.ParseConfig(): %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return DB{}, fmt.Errorf("create pool with config: %w", err)
	}

	return DB{
		Pool: pool,
	}, nil
}
