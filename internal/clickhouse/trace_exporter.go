package clickhouse

import (
	"context"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"go.opentelemetry.io/otel/sdk/trace"
)

type TraceExporter struct {
	serviceName string
	client      *ch.Client
}

func NewTraceExporter(serviceName string, client *ch.Client) *TraceExporter {
	return &TraceExporter{
		serviceName: serviceName,
		client:      client,
	}
}

func (e *TraceExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	// initialize columns matching otel_traces table schema
	colTimestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMicro)
	colTraceId := new(proto.ColStr)
	colSpanId := new(proto.ColStr)
	colParentSpanId := new(proto.ColStr)
	colSpanName := new(proto.ColStr).LowCardinality()
	colSpanKind := new(proto.ColStr).LowCardinality()
	colServiceName := new(proto.ColStr).LowCardinality()
	colResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colSpanAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colDuration := new(proto.ColUInt64)
	colStatusCode := new(proto.ColStr).LowCardinality()
	colStatusMessage := new(proto.ColStr)

	// map otel data to columns
	for _, span := range spans {
		colTimestamp.Append(span.StartTime())
		colTraceId.Append(span.SpanContext().TraceID().String())
		colSpanId.Append(span.SpanContext().SpanID().String())
		colParentSpanId.Append(span.Parent().SpanID().String())
		colSpanName.Append(span.Name())
		colSpanKind.Append(span.SpanKind().String())
		colServiceName.Append(e.serviceName)

		resourceAttributes := make(map[string]string)
		for _, attr := range span.Resource().Attributes() {
			resourceAttributes[string(attr.Key)] = attr.Value.AsString()
		}
		colResourceAttributes.Append(resourceAttributes)

		spanAttributes := make(map[string]string)
		for _, attr := range span.Attributes() {
			spanAttributes[string(attr.Key)] = attr.Value.AsString()
		}
		colSpanAttributes.Append(spanAttributes)

		colDuration.Append(uint64(span.EndTime().Sub(span.StartTime())))
		colStatusCode.Append(span.Status().Code.String())
		colStatusMessage.Append(span.Status().Description)
	}

	// construct input block
	input := proto.Input{
		{Name: "Timestamp", Data: colTimestamp},
		{Name: "TraceId", Data: colTraceId},
		{Name: "SpanId", Data: colSpanId},
		{Name: "ParentSpanId", Data: colParentSpanId},
		{Name: "SpanName", Data: colSpanName},
		{Name: "SpanKind", Data: colSpanKind},
		{Name: "ServiceName", Data: colServiceName},
		{Name: "ResourceAttributes", Data: colResourceAttributes},
		{Name: "SpanAttributes", Data: colSpanAttributes},
		{Name: "Duration", Data: colDuration},
		{Name: "StatusCode", Data: colStatusCode},
		{Name: "StatusMessage", Data: colStatusMessage},
	}

	// execute native protocol insert
	err := e.client.Do(ctx, ch.Query{
		Body:  "INSERT INTO otel_traces VALUES",
		Input: input,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *TraceExporter) Shutdown(ctx context.Context) error {
	if e.client != nil {
		return e.client.Close()
	}

	return nil
}
