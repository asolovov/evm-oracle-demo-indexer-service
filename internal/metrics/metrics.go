// Package metrics holds the indexer-service Prometheus registry and
// the small adapter that satisfies the chainsub + stream-hub metrics
// interfaces.
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/streamhub"
)

// Registry wraps prometheus.Registry so /metrics serves from a
// service-owned collector set (no global state leaks). The surface is
// deliberately small after the confirmation gate was removed — there
// is no confirmed/orphaned/reorg lifecycle any more, only "seen".
type Registry struct {
	r *prometheus.Registry

	EventsSeen   *prometheus.CounterVec // indexer_events_seen_total{kind}
	StreamDrops  *prometheus.CounterVec // indexer_stream_drops_total{reason}
	DecodeErrors prometheus.Counter     // indexer_decode_errors_total
}

// New constructs the metrics surface.
func New() *Registry {
	reg := prometheus.NewRegistry()
	r := &Registry{
		r: reg,
		EventsSeen: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "indexer_events_seen_total",
				Help: "Events persisted-and-published, by kind (emit-on-ingest; no confirmation gate).",
			},
			[]string{"kind"},
		),
		StreamDrops: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "indexer_stream_drops_total",
				Help: "StreamEvents subscriber drops, by reason.",
			},
			[]string{"reason"},
		),
		DecodeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "indexer_decode_errors_total",
			Help: "Logs the chainsub parser failed to decode.",
		}),
	}
	reg.MustRegister(r.EventsSeen, r.StreamDrops, r.DecodeErrors)
	return r
}

// Handler serves the registry on /metrics.
func (r *Registry) Handler() http.Handler {
	return promhttp.HandlerFor(r.r, promhttp.HandlerOpts{Registry: r.r})
}

// Chainsub satisfies chainsub.Metrics.
type Chainsub struct{ R *Registry }

// ObserveSeen increments indexer_events_seen_total{kind}.
func (c Chainsub) ObserveSeen(k models.EventKind) {
	c.R.EventsSeen.WithLabelValues(k.String()).Inc()
}

// ObserveDecodeError increments indexer_decode_errors_total.
func (c Chainsub) ObserveDecodeError() { c.R.DecodeErrors.Inc() }

// HubDrop returns a streamhub.DropFunc that bumps indexer_stream_drops_total.
func (r *Registry) HubDrop() streamhub.DropFunc {
	return func(_ uint64, reason streamhub.DropReason) {
		r.StreamDrops.WithLabelValues(string(reason)).Inc()
	}
}
