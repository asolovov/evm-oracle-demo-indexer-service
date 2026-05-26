// Package healthz hosts the indexer-service's tiny HTTP listener for
// liveness (/healthz), readiness (/readyz), and metrics (/metrics).
// Keeping these on a separate port from the gRPC server means
// orchestrators can probe without holding open gRPC connections.
package healthz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asolovov/evm-oracle-demo-indexer-service/config"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/module"
)

// ReadyChecker reports application readiness. Returns a non-nil
// error when any registered component is failing. In production this
// is satisfied by a thin adapter around module.Manager.HealthCheckAll.
type ReadyChecker interface {
	Ready(ctx context.Context) error
}

// Module wires the listener into the App's module.Manager lifecycle.
type Module struct {
	cfg          *config.HealthzConfig
	checker      ReadyChecker
	metricsH     http.Handler
	authorMeta   map[string]string

	srv *http.Server
}

// NewModule constructs an unstarted healthz module.
// `metricsHandler` may be nil — /metrics returns a placeholder if so.
func NewModule(cfg *config.HealthzConfig, checker ReadyChecker, metricsHandler http.Handler, authorMeta map[string]string) *Module {
	return &Module{
		cfg:        cfg,
		checker:    checker,
		metricsH:   metricsHandler,
		authorMeta: authorMeta,
	}
}

// Name implements module.Module.
func (m *Module) Name() string { return "healthz" }

// Init binds the listener but doesn't yet accept connections.
func (m *Module) Init(_ context.Context) error {
	if m.cfg == nil {
		return errors.New("healthz config is required")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", m.healthz)
	mux.HandleFunc("/readyz", m.readyz)

	if m.metricsH != nil {
		mux.Handle("/metrics", m.metricsH)
	} else {
		mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "metrics handler not wired", http.StatusNotImplemented)
		})
	}

	m.srv = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	return nil
}

// Start serves in a background goroutine.
func (m *Module) Start(_ context.Context) error {
	go func() {
		if err := m.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Surface via the application logger at the call site;
			// quiet here to keep ListenAndServe's normal shutdown from
			// looking like an error.
			_ = err
		}
	}()
	return nil
}

// Stop drains the listener.
func (m *Module) Stop(ctx context.Context) error {
	if m.srv == nil {
		return nil
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return m.srv.Shutdown(shutdownCtx)
}

// HealthCheck is a self-probe — always succeeds once Init has bound
// the mux. Module-graph readiness is what /readyz reports.
func (m *Module) HealthCheck(_ context.Context) error {
	if m.srv == nil {
		return errors.New("healthz module not initialised")
	}
	return nil
}

// healthz is the liveness probe — 200 once the listener is up.
func (m *Module) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{"status": "ok"}
	for k, v := range m.authorMeta {
		resp[k] = v
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// readyz walks every module's HealthCheck and returns 503 on any failure.
func (m *Module) readyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if m.checker == nil {
		http.Error(w, `{"status":"no_checker"}`, http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	if err := m.checker.Ready(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "not_ready", "error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// Ensure module.Module is satisfied at compile time.
var _ module.Module = (*Module)(nil)
