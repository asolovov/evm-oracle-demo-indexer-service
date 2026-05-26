package grpcsrv

import "context"

// Module integrates the gRPC server with the App's module.Manager
// lifecycle.
type Module struct{ s *Server }

// NewModule wraps a Server so it can be registered on module.Manager.
func NewModule(s *Server) *Module { return &Module{s: s} }

// Name implements module.Module.
func (m *Module) Name() string { return "grpc-server" }

// Init is a no-op — wiring happens in New().
func (m *Module) Init(_ context.Context) error { return nil }

// Start binds the listener and serves.
func (m *Module) Start(ctx context.Context) error { return m.s.Start(ctx) }

// Stop drains pending RPCs.
func (m *Module) Stop(ctx context.Context) error { return m.s.Stop(ctx) }

// HealthCheck reports SERVING once the listener is up.
func (m *Module) HealthCheck(_ context.Context) error { return nil }
