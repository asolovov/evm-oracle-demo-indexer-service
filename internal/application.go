// Package internal contains the indexer-service application wiring.
//
// Per architecture rules 1 + 2, this is the SINGLE place that
// constructs and connects modules. cmd/ does CLI + config load only;
// every component is built here and dependencies are passed in
// explicitly (no global state, no module self-wiring).
package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/asolovov/evm-oracle-demo-indexer-service/config"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/chainsub"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/grpcsrv"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/healthz"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/metrics"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/module"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/repository"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/streamhub"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/logger"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/version"
)

// App is the indexer-service application.
type App struct {
	config  *config.Scheme
	version *version.Version
	modules *module.Manager

	// chainsub is the one background goroutine the App owns directly
	// (it is not a module.Module — it has its own connect/reconnect
	// lifecycle). The hub is shared between chainsub (publisher) and
	// the gRPC server (subscriber source).
	chainsub *chainsub.Subscriber
	hub      *streamhub.Hub

	cancelRuntime context.CancelFunc
	runtimeWG     sync.WaitGroup
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

// Init constructs and registers everything:
//
//  1. Validate config (fail fast, architecture rule 6).
//  2. Repository module (opens the pgx pool).
//  3. Stream hub (shared) + chainsub (owns the chain client in one
//     goroutine; catches up then live-subscribes; emits on ingest).
//  4. gRPC server module + healthz module.
func (app *App) Init() error {
	if err := config.Validate(app.config); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	mts := metrics.New()

	// 1. Repository.
	repoModule := repository.NewModule(app.config.Database)
	app.modules.Register(repoModule)
	if err := repoModule.Init(context.Background()); err != nil {
		return fmt.Errorf("init repository: %w", err)
	}
	repo := repoModule.Repository()

	// 2. Stream hub — fed by chainsub, consumed by the gRPC server.
	app.hub = streamhub.New(app.config.Indexer.StreamSubscriberBuffer, mts.HubDrop())

	// 3. Chain subscriber. Owns its client end-to-end (no sharing →
	//    no data race). Publishes each persisted event to the hub.
	app.chainsub = chainsub.New(repo, app.hub, chainsub.Config{
		WSURL:           app.config.Chain.WSURL,
		RPCURL:          app.config.Chain.RPCURL,
		RegistryAddress: common.HexToAddress(app.config.Chain.RegistryAddress),
		DefaultStart:    app.config.Chain.BackfillFromBlock,
		ChunkSize:       app.config.Indexer.BackfillChunkSize,
		Metrics:         metrics.Chainsub{R: mts},
	})

	// 4. gRPC server + healthz.
	srv := grpcsrv.New(app.config.GRPC, repo, app.hub)
	app.modules.Register(grpcsrv.NewModule(srv))

	authorMeta := map[string]string{
		"version": app.version.String(),
		"name":    "evm-oracle-demo-indexer-service",
		"chain":   app.config.Chain.Name,
	}
	hz := healthz.NewModule(app.config.Healthz, &readyAdapter{m: app.modules}, mts.Handler(), authorMeta)
	app.modules.Register(hz)
	if err := hz.Init(context.Background()); err != nil {
		return fmt.Errorf("init healthz: %w", err)
	}

	return nil
}

// Serve launches the chainsub goroutine and the transport modules, then
// blocks until a termination signal arrives.
func (app *App) Serve() error {
	//nolint:gosec // cancel is stored in app.cancelRuntime and invoked from app.Stop; deferring here would kill the runtime the moment Serve returns.
	runtimeCtx, cancel := context.WithCancel(context.Background())
	app.cancelRuntime = cancel

	// chainsub owns its own connect/catch-up/live/reconnect lifecycle.
	app.runtimeWG.Add(1)
	go func() {
		defer app.runtimeWG.Done()
		if err := app.chainsub.Run(runtimeCtx); err != nil {
			logger.Log().Errorf("chainsub.Run: %v", err)
		}
	}()

	// Transport: gRPC server + healthz.
	if err := app.modules.StartAll(runtimeCtx); err != nil {
		return fmt.Errorf("start modules: %w", err)
	}

	logger.Log().Infof("indexer-service is running on chain=%s (chain_id=%d); ctrl-c to stop",
		app.config.Chain.Name, app.config.Chain.ChainID)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logger.Log().Info("shutdown signal received, stopping gracefully...")
	return nil
}

// Stop cancels the runtime goroutine and drains modules. Idempotent.
func (app *App) Stop() error {
	if app.cancelRuntime != nil {
		app.cancelRuntime()
	}
	if app.hub != nil {
		app.hub.Shutdown()
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	moduleErr := app.modules.StopAll(stopCtx)

	doneCh := make(chan struct{})
	go func() {
		app.runtimeWG.Wait()
		close(doneCh)
	}()
	select {
	case <-doneCh:
	case <-stopCtx.Done():
		logger.Log().Warn("runtime goroutine did not exit within 30s")
	}
	return moduleErr
}

// Config returns the application configuration.
func (app *App) Config() *config.Scheme { return app.config }

// Version returns the application version string.
func (app *App) Version() string { return app.version.String() }

// Modules returns the module manager.
func (app *App) Modules() *module.Manager { return app.modules }

// readyAdapter translates module.Manager.HealthCheckAll's map into a
// single non-nil error for healthz.ReadyChecker.
type readyAdapter struct{ m *module.Manager }

func (a *readyAdapter) Ready(ctx context.Context) error {
	failures := a.m.HealthCheckAll(ctx)
	var bad []string
	for name, err := range failures {
		if err != nil {
			bad = append(bad, fmt.Sprintf("%s: %v", name, err))
		}
	}
	if len(bad) == 0 {
		return nil
	}
	return errors.New("modules unhealthy: " + strings.Join(bad, "; "))
}
