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

	"github.com/ClickHouse/ch-go"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type System struct {
	cfg    *config.AppConfig
	db     *sql.DB
	ch     *ch.Client
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger *slog.Logger
	mux    *chi.Mux
	rpc    *grpc.Server
	waiter waiter.Waiter
}

func NewSystem(ctx context.Context, cfg *config.AppConfig) (*System, error) {
	// connect postgres database
	pgConn, err := sql.Open("postgres", cfg.PG.Conn)
	if err != nil {
		return nil, fmt.Errorf("pg connection: %w", err)
	}

	// connect clickhouse database
	chConn, err := ch.Dial(ctx, ch.Options{
		Address:  cfg.CH.Address(),
		Database: cfg.CH.Database,
		User:     cfg.CH.Username,
		Password: cfg.CH.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("ch connection: %w", err)
	}

	// connect nats jetstream
	ntConn, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection: %w", err)
	}

	jsc, err := ntConn.JetStream()
	if err != nil {
		return nil, err
	}

	// init jetstream
	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     cfg.Nats.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", cfg.Nats.Stream)},
	})
	if err != nil {
		return nil, err
	}

	// // init otel propagation
	// textMapPropagator := propagation.NewCompositeTextMapPropagator(
	// 	propagation.TraceContext{},
	// 	propagation.Baggage{},
	// )
	// otel.SetTextMapPropagator(textMapPropagator)

	// // init otel traces
	// tracerProvider := trace.NewTracerProvider(trace.WithBatcher(
	// 	clickhouse.NewTraceExporter(cfg.ServiceName, chConn),
	// 	trace.WithBatchTimeout(time.Second),
	// ))
	// otel.SetTracerProvider(tracerProvider)

	// // init otel metrics
	// meterProvider := metric.NewMeterProvider(metric.WithReader(
	// 	metric.NewPeriodicReader(
	// 		clickhouse.NewMetricExporter(cfg.ServiceName, chConn),
	// 		metric.WithInterval(3*time.Second),
	// 	),
	// ))
	// otel.SetMeterProvider(meterProvider)

	// // init otel logs
	// loggerProvider := log.NewLoggerProvider(log.WithProcessor(
	// 	log.NewBatchProcessor(
	// 		clickhouse.NewLogExpoerter(cfg.ServiceName, chConn),
	// 	),
	// ))
	// global.SetLoggerProvider(loggerProvider)

	// init logger
	logger := logger.New(logger.LogConfig{
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})

	// init grpc
	grpcServer := grpc.NewServer( /*grpc.StatsHandler(otelgrpc.NewServerHandler())*/ )
	reflection.Register(grpcServer)

	// init infrastructure
	s := &System{
		cfg:    cfg,
		db:     pgConn,
		ch:     chConn,
		nc:     ntConn,
		js:     jsc,
		logger: logger,
		mux:    chi.NewMux(),
		rpc:    grpcServer,
		waiter: waiter.New(waiter.CatchSignals()),
	}

	// // cleanup otel providers
	// s.waiter.AddCleanupFunc(func() {
	// 	if err := tracerProvider.Shutdown(ctx); err != nil {
	// 		s.logger.ErrorContext(ctx, "error shutting down the tracer provider", "error", err.Error())
	// 	}
	// })

	// s.waiter.AddCleanupFunc(func() {
	// 	if err := meterProvider.Shutdown(ctx); err != nil {
	// 		logger.ErrorContext(ctx, "error shutting down the metric provider", "error", err.Error())
	// 	}
	// })

	// s.waiter.AddCleanupFunc(func() {
	// 	if err := loggerProvider.Shutdown(ctx); err != nil {
	// 		logger.ErrorContext(ctx, "error shutting down the logger provider", "error", err.Error())
	// 	}
	// })

	return s, nil
}

func (s *System) Config() *config.AppConfig { return s.cfg }
func (s *System) DB() *sql.DB               { return s.db }
func (s *System) CH() *ch.Client            { return s.ch }
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
		s.logger.InfoContext(ctx, "web server started", "address", s.cfg.Web.Address())
		defer s.logger.InfoContext(ctx, "web server shutdown")

		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.InfoContext(ctx, "web server to be shutdown")

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
		s.logger.InfoContext(ctx, "rpc server started", "address", s.cfg.Rpc.Address())
		defer s.logger.InfoContext(ctx, "rpc server shutdown")

		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.logger.InfoContext(ctx, "rpc server to be shutdown")

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
		s.logger.InfoContext(ctx, "message stream started")
		defer s.logger.InfoContext(ctx, "message stream stopped")
		<-closed
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		return s.nc.Drain()
	})

	return group.Wait()
}
