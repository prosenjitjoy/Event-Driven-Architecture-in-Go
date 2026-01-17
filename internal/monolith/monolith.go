package monolith

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/internal/config"
	"mall/internal/waiter"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type Monolith interface {
	Config() *config.AppConfig
	DB() *sql.DB
	JS() nats.JetStreamContext
	Logger() *slog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
}

type Module interface {
	Startup(context.Context, Monolith) error
}
