package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/asolovov/evm-oracle-demo-indexer-service/config"
)

// Module integrates the Repository into the App's module.Manager
// lifecycle so app start/stop cleanly opens and drains the pool.
type Module struct {
	cfg  *config.DatabaseConfig
	pool *pgxpool.Pool
	repo *Repository
}

// NewModule creates an unstarted Module.
func NewModule(cfg *config.DatabaseConfig) *Module { return &Module{cfg: cfg} }

// Name implements module.Module.
func (m *Module) Name() string { return "repository" }

// Init opens the pgx pool and pings it.
func (m *Module) Init(ctx context.Context) error {
	if m.cfg == nil {
		return fmt.Errorf("repository module: database config is required")
	}
	// Parse a DSN that does NOT contain the password, then set the
	// password as a struct field. A ParseConfig failure echoes the
	// DSN it was given into the error; keeping the secret out of that
	// string means it can never leak into logs (org rule 1).
	cfg, err := pgxpool.ParseConfig(buildDSN(m.cfg))
	if err != nil {
		return fmt.Errorf("parse db config: %w", err)
	}
	cfg.ConnConfig.Password = m.cfg.Password
	if m.cfg.MaxOpenConns > 0 {
		cfg.MaxConns = int32(m.cfg.MaxOpenConns) //nolint:gosec // small bounded value.
	}
	if m.cfg.MaxIdleConns > 0 {
		cfg.MinConns = int32(m.cfg.MaxIdleConns) //nolint:gosec // small bounded value.
	}
	if m.cfg.ConnMaxLifetime > 0 {
		cfg.MaxConnLifetime = time.Duration(m.cfg.ConnMaxLifetime) * time.Second
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to evm_indexer: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("ping evm_indexer: %w", err)
	}
	m.pool = pool
	m.repo = New(pool)
	return nil
}

// Start is a no-op — the pool is already active after Init.
func (m *Module) Start(_ context.Context) error { return nil }

// Stop closes the pool.
func (m *Module) Stop(_ context.Context) error {
	if m.pool != nil {
		m.pool.Close()
	}
	return nil
}

// HealthCheck pings the pool.
func (m *Module) HealthCheck(ctx context.Context) error {
	if m.pool == nil {
		return fmt.Errorf("repository module not initialized")
	}
	return m.pool.Ping(ctx)
}

// Repository exposes the typed repository for downstream modules.
func (m *Module) Repository() *Repository { return m.repo }

// buildDSN omits the password on purpose — it is set on
// cfg.ConnConfig.Password after parsing so it can never appear in a
// ParseConfig error string (org rule 1: no creds in output).
func buildDSN(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)
}
