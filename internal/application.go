// Package internal contains the indexer-service application wiring.
package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asolovov/evm-oracle-demo-indexer-service/config"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/module"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/logger"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/version"
)

// App is the indexer-service application instance. The full module
// graph is wired in registerModules — additions are added in the order
// db pool -> repos -> chain client -> stream hub -> confirmer -> chainsub
// -> backfill -> gRPC server -> healthz, with each later stage receiving
// dependencies as constructor arguments (architecture rules 1+2).
type App struct {
	config  *config.Scheme
	version *version.Version
	modules *module.Manager
}

// NewApplication creates a new App instance.
func NewApplication() (*App, error) {
	ver, err := version.NewVersion()
	if err != nil {
		return nil, fmt.Errorf("init app version: %w", err)
	}

	return &App{
		config:  &config.Scheme{},
		version: ver,
		modules: module.NewManager(),
	}, nil
}

// Init initializes the application and all registered modules.
func (app *App) Init() error {
	if err := app.registerModules(); err != nil {
		return fmt.Errorf("register modules: %w", err)
	}
	return nil
}

// registerModules wires the indexer-service module graph. Filled in by
// later commits — for now this is a no-op so the renamed template
// scaffolding compiles cleanly.
func (app *App) registerModules() error {
	logger.Log().Infof("registered %d modules (indexer-service skeleton)", app.modules.Count())
	return nil
}

// Serve starts all modules and waits for shutdown signal.
func (app *App) Serve() error {
	ctx := context.Background()

	if err := app.modules.StartAll(ctx); err != nil {
		return fmt.Errorf("start modules: %w", err)
	}

	logger.Log().Info("indexer-service is running; press Ctrl+C to stop")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logger.Log().Info("shutdown signal received, stopping gracefully...")
	return nil
}

// Stop gracefully shuts down all modules.
func (app *App) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return app.modules.StopAll(ctx)
}

// Config returns the application configuration.
func (app *App) Config() *config.Scheme { return app.config }

// Version returns the application version string.
func (app *App) Version() string { return app.version.String() }

// Modules returns the module manager (useful for health checks).
func (app *App) Modules() *module.Manager { return app.modules }
