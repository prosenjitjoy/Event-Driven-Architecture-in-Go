package postgresotel

import (
	"context"
	"database/sql"
	"errors"
	"mall/internal/postgres"
	"time"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedDB struct {
	db postgres.DBTX
}

var _ postgres.DBTX = (*tracedDB)(nil)

func Trace(db postgres.DBTX) postgres.DBTX {
	return tracedDB{db: db}
}

func (t tracedDB) PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("PrepareContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.db.PrepareContext(ctx, query)
}

func (t tracedDB) ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("ExecContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.db.ExecContext(ctx, query, args...)
}

func (t tracedDB) QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("QueryContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.db.QueryContext(ctx, query, args...)
}

func (t tracedDB) QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("QueryRowContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", time.Since(started).Seconds()),
		))
		t.recordError(span, row.Err())
	}(time.Now())

	return t.db.QueryRowContext(ctx, query, args...)
}

func (t tracedDB) recordError(span trace.Span, err error) {
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			span.AddEvent("Database Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
				attribute.String("Code", string(pqErr.Code)),
				attribute.String("Severity", pqErr.Severity),
				attribute.String("Message", pqErr.Message),
				attribute.String("Detail", pqErr.Detail),
			))
		} else {
			span.AddEvent("Database Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
			))
		}
	}
}
