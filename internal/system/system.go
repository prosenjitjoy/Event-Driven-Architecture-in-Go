package system

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"mall/internal/config"
	"mall/internal/logger"
	"mall/internal/waiter"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type System struct {
	cfg    *config.AppConfig
	db     *sql.DB
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger *slog.Logger
	mux    *chi.Mux
	rpc    *grpc.Server
	waiter waiter.Waiter
}

func NewSystem(cfg *config.AppConfig) (*System, error) {
	// connect database
	db, err := sql.Open("postgres", cfg.PG.Conn)
	if err != nil {
		return nil, fmt.Errorf("db connection: %w", err)
	}

	// connect nats jetstream
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// init jetstream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.Nats.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", cfg.Nats.Stream)},
	})
	if err != nil {
		return nil, err
	}

	// init logger
	logger := logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})

	// init grpc
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// init infrastructure
	s := &System{
		cfg:    cfg,
		db:     db,
		nc:     nc,
		js:     js,
		logger: logger,
		mux:    chi.NewMux(),
		rpc:    grpcServer,
		waiter: waiter.New(waiter.CatchSignals()),
	}

	return s, nil
}

func (s *System) Config() *config.AppConfig { return s.cfg }
func (s *System) DB() *sql.DB               { return s.db }
func (s *System) NC() *nats.Conn            { return s.nc }
func (s *System) JS() nats.JetStreamContext { return s.js }
func (s *System) Logger() *slog.Logger      { return s.logger }
func (s *System) Mux() *chi.Mux             { return s.mux }
func (s *System) RPC() *grpc.Server         { return s.rpc }
func (s *System) Waiter() waiter.Waiter     { return s.waiter }

func (s *System) MigrateDB(fs fs.FS) error {
	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(s.db, "."); err != nil {
		return err
	}

	return nil
}

func (s *System) WaitForWeb(ctx context.Context) error {
	webServer := &http.Server{
		Addr:    s.cfg.Web.Address(),
		Handler: s.mux,
	}

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		s.logger.Info("web server started", "address", s.cfg.Web.Address())
		defer s.logger.Info("web server shutdown")

		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.Info("web server to be shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}

func (s *System) WaitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.Rpc.Address())
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		s.logger.Info("rpc server started", "address", s.cfg.Rpc.Address())
		defer s.logger.Info("rpc server shutdown")

		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.Info("rpc server to be shutdown")

		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()

		timeout := time.NewTimer(s.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			s.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (s *System) WaitForStream(ctx context.Context) error {
	closed := make(chan struct{})
	s.nc.SetClosedHandler(func(c *nats.Conn) {
		close(closed)
	})

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		s.logger.Info("message stream started")
		defer s.logger.Info("message stream stopped")
		<-closed
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		return s.nc.Drain()
	})

	return group.Wait()
}
