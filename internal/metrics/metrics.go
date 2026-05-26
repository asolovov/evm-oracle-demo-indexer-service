// Package metrics holds the indexer-service Prometheus registry and
// the small adapter types that satisfy the metrics interfaces declared
// by chainsub, confirmer, and the stream hub.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/streamhub"
)

// Registry wraps prometheus.Registry so /metrics can serve from a
// service-owned collector set (no global state leaks).
type Registry struct {
	r *prometheus.Registry

	EventsTotal           *prometheus.CounterVec
	LagSeconds            prometheus.Gauge
	ReorgTotal            prometheus.Counter
	StreamSubscribers     *prometheus.GaugeVec
	StreamDrops           *prometheus.CounterVec
	DecodeErrors          prometheus.Counter
}

// New constructs the metrics surface.
func New() *Registry {
	reg := prometheus.NewRegistry()

	r := &Registry{
		r: reg,

		EventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "indexer_events_total",
				Help: "Lifecycle counters per event kind. status is one of seen|confirmed|orphaned.",
			},
			[]string{"kind", "status"},
		),
		LagSeconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "indexer_lag_seconds",
			Help: "Seconds between the chain head and the last block the indexer has fully processed (best-effort wall-time estimate).",
		}),
		ReorgTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "indexer_reorg_total",
			Help: "Cumulative count of events the confirmer marked orphaned due to a reorg.",
		}),
		StreamSubscribers: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "indexer_stream_subscribers",
				Help: "Active StreamEvents subscribers, labelled by kind filter.",
			},
			[]string{"kind"},
		),
		StreamDrops: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "indexer_stream_drops_total",
				Help: "StreamEvents subscriber drops, labelled by reason.",
			},
			[]string{"reason"},
		),
		DecodeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "indexer_decode_errors_total",
			Help: "Cumulative count of logs the chainsub parser failed to decode.",
		}),
	}

	reg.MustRegister(
		r.EventsTotal,
		r.LagSeconds,
		r.ReorgTotal,
		r.StreamSubscribers,
		r.StreamDrops,
		r.DecodeErrors,
	)
	return r
}

// Handler returns an http.Handler that serves the registry on /metrics.
func (r *Registry) Handler() http.Handler {
	return promhttp.HandlerFor(r.r, promhttp.HandlerOpts{Registry: r.r})
}

// Chainsub satisfies chainsub.Metrics.
type Chainsub struct{ R *Registry }

// ObserveSeen increments indexer_events_total{kind,status=seen}.
func (c Chainsub) ObserveSeen(k models.EventKind) {
	c.R.EventsTotal.WithLabelValues(k.String(), "seen").Inc()
}

// ObserveDecodeError increments indexer_decode_errors_total.
func (c Chainsub) ObserveDecodeError() { c.R.DecodeErrors.Inc() }

// Confirmer satisfies confirmer.Metrics.
type Confirmer struct{ R *Registry }

// ObserveConfirmed increments indexer_events_total{kind,status=confirmed}.
func (c Confirmer) ObserveConfirmed(k models.EventKind) {
	c.R.EventsTotal.WithLabelValues(k.String(), "confirmed").Inc()
}

// ObserveOrphaned increments indexer_events_total{kind,status=orphaned} + indexer_reorg_total.
func (c Confirmer) ObserveOrphaned(k models.EventKind) {
	c.R.EventsTotal.WithLabelValues(k.String(), "orphaned").Inc()
	c.R.ReorgTotal.Inc()
}

// ObserveLagSeconds sets indexer_lag_seconds.
func (c Confirmer) ObserveLagSeconds(v float64) { c.R.LagSeconds.Set(v) }

// HubDrop returns a streamhub.DropFunc that bumps indexer_stream_drops_total.
func (r *Registry) HubDrop() streamhub.DropFunc {
	return func(_ uint64, reason streamhub.DropReason) {
		r.StreamDrops.WithLabelValues(string(reason)).Inc()
	}
}
