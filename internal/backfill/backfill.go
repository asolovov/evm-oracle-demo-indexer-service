// Package backfill closes the gap between the persistent
// `chain_cursor` and the live WS subscription. On startup, before
// (or in parallel with) the chainsub goroutine, it walks
// `[cursor+1, head]` in fixed-size chunks via `eth_getLogs`, decodes,
// and persists via the same parser/repo path the live subscription
// uses — so confirmer + stream hub see a single, ordered firehose.
package backfill

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/chainsub"
	"github.com/asolovov/evm-oracle-demo-indexer-service/internal/models"
	"github.com/asolovov/evm-oracle-demo-indexer-service/pkg/logger"
)

// LogFetcher is the chain surface backfill needs. Matches
// `*ethclient.Client` in production.
type LogFetcher interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error)
}

// Store is the persistence surface backfill needs. Same shape as the
// chainsub one — InsertEvent is replay-idempotent.
type Store interface {
	InsertEvent(ctx context.Context, e *models.Event) (bool, error)
	ChainCursor(ctx context.Context) (*models.ChainCursor, error)
	UpdateChainCursor(ctx context.Context, block uint64) error
	AggregatorRegistry(ctx context.Context) (map[common.Address]common.Hash, error)
}

// Metrics is optional; nil-safe.
type Metrics interface {
	ObserveSeen(kind models.EventKind)
	ObserveDecodeError()
}

// Reconciler walks the gap [cursor+1, head] in chunks.
type Reconciler struct {
	client        LogFetcher
	store         Store
	parser        *chainsub.Parser
	registryAddr  common.Address
	defaultStart  uint64
	chunkSize     uint64
	cursorEvery   uint64
	metrics       Metrics
}

// Config bundles the constructor knobs.
type Config struct {
	RegistryAddress common.Address
	DefaultStart    uint64 // used when chain_cursor.last_processed_block is 0 / unset.
	ChunkSize       uint64
	CursorEvery     uint64 // persist cursor every N blocks (clamped to >= 1).
	Metrics         Metrics
}

// New builds a Reconciler. `parser` must be configured with the same
// aggregator resolver the live subscriber uses so PriceRequested/
// Fulfilled events get the correct asset_id from address.
func New(client LogFetcher, store Store, parser *chainsub.Parser, cfg Config) *Reconciler {
	if cfg.ChunkSize == 0 {
		cfg.ChunkSize = 1000
	}
	if cfg.CursorEvery == 0 {
		cfg.CursorEvery = 100
	}
	return &Reconciler{
		client:       client,
		store:        store,
		parser:       parser,
		registryAddr: cfg.RegistryAddress,
		defaultStart: cfg.DefaultStart,
		chunkSize:    cfg.ChunkSize,
		cursorEvery:  cfg.CursorEvery,
		metrics:      cfg.Metrics,
	}
}

// Run walks the gap. Returns nil on success or ctx cancel; otherwise
// the error from the underlying call. Run is single-shot — callers
// invoke it during App.Start before (or alongside) the live
// subscription.
func (r *Reconciler) Run(ctx context.Context) error {
	cursor, err := r.store.ChainCursor(ctx)
	if err != nil {
		return fmt.Errorf("read chain cursor: %w", err)
	}
	startFrom := cursor.LastProcessedBlock + 1
	if cursor.LastProcessedBlock == 0 && r.defaultStart > 0 {
		startFrom = r.defaultStart
	}

	head, err := r.client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("head lookup: %w", err)
	}
	headNum := head.Number.Uint64()
	if startFrom > headNum {
		logger.Log().Infof("backfill: nothing to do — cursor=%d head=%d", cursor.LastProcessedBlock, headNum)
		return nil
	}

	addresses, err := r.addresses(ctx)
	if err != nil {
		return err
	}

	logger.Log().Infof("backfill: replaying [%d..%d] in %d-block chunks (%d address(es))",
		startFrom, headNum, r.chunkSize, len(addresses))

	from := startFrom
	persistAt := startFrom + r.cursorEvery - 1
	started := time.Now()

	for from <= headNum {
		if err := ctx.Err(); err != nil {
			return err
		}
		to := from + r.chunkSize - 1
		if to > headNum {
			to = headNum
		}

		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from),
			ToBlock:   new(big.Int).SetUint64(to),
			Addresses: addresses,
		}

		logs, err := r.client.FilterLogs(ctx, query)
		if err != nil {
			return fmt.Errorf("eth_getLogs [%d..%d]: %w", from, to, err)
		}
		for _, log := range logs {
			if perr := r.handleLog(ctx, log); perr != nil {
				logger.Log().Warnf("backfill: handleLog [%d..%d]: %v", from, to, perr)
			}
		}

		if to >= persistAt {
			if err := r.store.UpdateChainCursor(ctx, to); err != nil {
				return fmt.Errorf("persist cursor at %d: %w", to, err)
			}
			persistAt = to + r.cursorEvery
		}
		from = to + 1
	}

	// Final cursor flush so the next start picks up exactly where we left off.
	if err := r.store.UpdateChainCursor(ctx, headNum); err != nil {
		return fmt.Errorf("persist final cursor at %d: %w", headNum, err)
	}

	logger.Log().Infof("backfill: complete in %s — processed [%d..%d]", time.Since(started), startFrom, headNum)
	return nil
}

func (r *Reconciler) handleLog(ctx context.Context, log types.Log) error {
	evt, err := r.parser.Parse(log)
	if err != nil {
		if errors.Is(err, chainsub.ErrUnknownEvent) {
			return nil
		}
		if errors.Is(err, chainsub.ErrAssetIDUnknown) {
			// Still persist — the asset mapping might be filled in
			// later by a reconnect-time refresh.
			logger.Log().Warnf("backfill: aggregator %s unknown to resolver; persisting with empty asset_id", log.Address.Hex())
		} else {
			if r.metrics != nil {
				r.metrics.ObserveDecodeError()
			}
			return fmt.Errorf("parse log (tx=%s idx=%d): %w", log.TxHash.Hex(), log.Index, err)
		}
	}
	inserted, err := r.store.InsertEvent(ctx, evt)
	if err != nil {
		return fmt.Errorf("persist event: %w", err)
	}
	if inserted && r.metrics != nil {
		r.metrics.ObserveSeen(evt.Kind)
	}
	return nil
}

// addresses returns the registry + every known aggregator.
// Matches the chainsub subscription filter so events the live
// subscription would have picked up are also picked up on replay.
func (r *Reconciler) addresses(ctx context.Context) ([]common.Address, error) {
	mapping, err := r.store.AggregatorRegistry(ctx)
	if err != nil {
		return nil, fmt.Errorf("aggregator registry: %w", err)
	}
	out := make([]common.Address, 0, len(mapping)+1)
	out = append(out, r.registryAddr)
	for addr := range mapping {
		out = append(out, addr)
	}
	return out, nil
}
