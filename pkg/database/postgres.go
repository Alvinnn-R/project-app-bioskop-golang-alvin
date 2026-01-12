package database

import (
	"context"
	"fmt"
	"session-23/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(config utils.DatabaseCofig) (*pgxpool.Pool, error) {
	// connStr := "user=postgres password=postgres dbname=assignment_with_role sslmode=disable host=192.168.1.12"
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		config.Username, config.Password, config.Name, config.Host)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// ---- Pool settings ----
	cfg.MaxConns = config.MaxConn // max connection pool
	cfg.MinConns = 5              // min connection pool
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	// Optional: set timeout connect
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Test connection
	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return pool, nil
}
