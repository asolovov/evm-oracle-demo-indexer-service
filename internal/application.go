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
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/backfill"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/chainsub"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/confirmer"
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

	// Owned subsystems (kept on App so Stop can drain them in the
	// right order independently of the module.Manager).
	chainsub  *chainsub.Subscriber
	confirmer *confirmer.Confirmer
	backfill  *backfill.Reconciler
	hub       *streamhub.Hub

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

// Init constructs and registers every module. Splits cleanly into:
//
//  1. Validate config (fail fast, architecture rule 6).
//  2. Repository module (Init opens the pgx pool).
//  3. Build stream hub + confirmer + chainsub + backfill (plain
//     packages — these are not module.Module instances; their
//     goroutines are managed by the App's Run/Stop pair).
//  4. gRPC server module.
//  5. Healthz module — its readyz walks the module graph.
func (app *App) Init() error {
	if err := config.Validate(app.config); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	mts := metrics.New()

	// 1. Repository
	repoModule := repository.NewModule(app.config.Database)
	app.modules.Register(repoModule)
	if err := repoModule.Init(context.Background()); err != nil {
		return fmt.Errorf("init repository: %w", err)
	}
	repo := repoModule.Repository()

	// 2. Stream hub.
	app.hub = streamhub.New(app.config.Indexer.StreamSubscriberBuffer, mts.HubDrop())

	// 3. Chain subscriber (does NOT dial yet; Run() dials on first
	//    iteration). Confirmer + Backfill use the subscriber's
	//    *ethclient.Client after Run() opens it.
	regAddr := common.HexToAddress(app.config.Chain.RegistryAddress)
	app.chainsub = chainsub.New(repo, chainsub.Config{
		WSURL:           app.config.Chain.WSURL,
		RPCURL:          app.config.Chain.RPCURL,
		RegistryAddress: regAddr,
		ReconnectWait:   2 * time.Second,
		Metrics:         metrics.Chainsub{R: mts},
	})

	// 4. gRPC server.
	srv := grpcsrv.New(app.config.GRPC, repo, app.hub, app.config.Indexer.Confirmations)
	app.modules.Register(grpcsrv.NewModule(srv))

	// 5. Healthz.
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

	app.runtime(mts)
	return nil
}

// runtime captures the metrics surface for the background goroutines
// that Run() will launch in Serve().
func (app *App) runtime(mts *metrics.Registry) {
	app.confirmer = confirmer.New(nil, nil, app.hub, confirmer.Config{
		Threshold:  app.config.Indexer.Confirmations,
		Interval:   time.Duration(app.config.Indexer.ReorgCheckIntervalSec) * time.Second,
		BatchLimit: 500,
		Metrics:    metrics.Confirmer{R: mts},
	})
	// repo + chain client get injected at Serve time after the
	// chainsub goroutine dials.

	_ = mts // metrics registry is captured by every collector already.
}

// Serve starts everything in order. Returns when a signal terminates
// the process or the runtime goroutines exit.
//
// Lifecycle:
//
//   - chainsub.Run goroutine: dials WS+RPC, drains logs into the
//     repo, maintains the aggregator->asset mapping.
//   - backfill (one-shot): waits for chainsub to dial, then runs
//     the gap-fill pass against the same RPC client.
//   - confirmer.Run goroutine: ticks every reorg_check_interval_sec.
//   - module.Manager.StartAll: brings up the gRPC server + healthz.
func (app *App) Serve() error {
	runtimeCtx, cancel := context.WithCancel(context.Background())
	app.cancelRuntime = cancel

	// chainsub
	app.runtimeWG.Add(1)
	go func() {
		defer app.runtimeWG.Done()
		if err := app.chainsub.Run(runtimeCtx); err != nil {
			logger.Log().Errorf("chainsub.Run: %v", err)
		}
	}()

	// Wait briefly for the chainsub to dial so confirmer + backfill
	// have a client to share. The first iteration of chainsub.Run
	// must establish the client before backfill can borrow it; the
	// gating is best-effort, not perfectly synchronous.
	app.waitForChainClient(runtimeCtx, 30*time.Second)

	if app.chainsub.Client() == nil {
		logger.Log().Warn("chainsub never connected; backfill + confirmer running in degraded mode")
	}

	repoModule := app.findRepository()
	if repoModule == nil {
		return errors.New("repository module not registered")
	}
	repo := repoModule.Repository()

	// Backfill (one-shot).
	if app.chainsub.Client() != nil {
		parser, err := chainsub.NewParser(common.HexToAddress(app.config.Chain.RegistryAddress), app.chainsub)
		if err != nil {
			logger.Log().Warnf("backfill: parser init failed: %v — skipping", err)
		} else {
			app.backfill = backfill.New(app.chainsub.Client(), repo, parser, backfill.Config{
				RegistryAddress: common.HexToAddress(app.config.Chain.RegistryAddress),
				DefaultStart:    app.config.Chain.BackfillFromBlock,
				ChunkSize:       app.config.Indexer.BackfillChunkSize,
			})
			app.runtimeWG.Add(1)
			go func() {
				defer app.runtimeWG.Done()
				if err := app.backfill.Run(runtimeCtx); err != nil {
					logger.Log().Warnf("backfill: %v", err)
				}
			}()
		}
	}

	// Confirmer.
	if app.chainsub.Client() != nil {
		// Re-build confirmer with the now-available chain client + repo.
		app.confirmer = confirmer.New(repo, app.chainsub.Client(), app.hub, confirmer.Config{
			Threshold:  app.config.Indexer.Confirmations,
			Interval:   time.Duration(app.config.Indexer.ReorgCheckIntervalSec) * time.Second,
			BatchLimit: 500,
		})
		app.runtimeWG.Add(1)
		go func() {
			defer app.runtimeWG.Done()
			if err := app.confirmer.Run(runtimeCtx); err != nil {
				logger.Log().Errorf("confirmer.Run: %v", err)
			}
		}()
	}

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

// Stop cancels runtime goroutines and drains modules. Idempotent.
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
		logger.Log().Warn("runtime goroutines did not exit within 30s")
	}
	return moduleErr
}

// Config returns the application configuration.
func (app *App) Config() *config.Scheme { return app.config }

// Version returns the application version string.
func (app *App) Version() string { return app.version.String() }

// Modules returns the module manager.
func (app *App) Modules() *module.Manager { return app.modules }

// waitForChainClient polls for chainsub.Client() to become non-nil.
// Returns once the client is ready or the timeout elapses.
func (app *App) waitForChainClient(ctx context.Context, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if app.chainsub.Client() != nil {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
		}
	}
}

// findRepository walks the module manager for the repository module.
func (app *App) findRepository() *repository.Module {
	for _, m := range app.modules.List() {
		if rm, ok := m.(*repository.Module); ok {
			return rm
		}
	}
	return nil
}

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
