package system

import (
	"context"
	"database/sql"
	"log/slog"
	"mall/internal/config"
	"mall/internal/waiter"

	"github.com/ClickHouse/ch-go"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type Service interface {
	Config() *config.AppConfig
	DB() *sql.DB
	CH() *ch.Client
	NC() *nats.Conn
	JS() nats.JetStreamContext
	Logger() *slog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
}

type Module interface {
	Startup(context.Context, Service) error
}
