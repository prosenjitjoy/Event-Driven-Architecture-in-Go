package clickhouse

import (
	"context"
	"time"

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
	colTimestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colTraceId := new(proto.ColStr)
	colSpanId := new(proto.ColStr)
	colParentSpanId := new(proto.ColStr)
	colTraceState := new(proto.ColStr)
	colSpanName := new(proto.ColStr).LowCardinality()
	colSpanKind := new(proto.ColStr).LowCardinality()
	colServiceName := new(proto.ColStr).LowCardinality()
	colResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colScopeName := new(proto.ColStr)
	colScopeVersion := new(proto.ColStr)
	colSpanAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colDuration := new(proto.ColUInt64)
	colStatusCode := new(proto.ColStr).LowCardinality()
	colStatusMessage := new(proto.ColStr)
	colEventsTimestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano).Array()
	colEventsName := new(proto.ColStr).LowCardinality().Array()
	colEventsAttributes := proto.NewArray(proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr)))
	colLinksTraceId := new(proto.ColStr).Array()
	colLinksSpanId := new(proto.ColStr).Array()
	colLinksTraceState := new(proto.ColStr).Array()
	colLinksAttributes := proto.NewArray(proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr)))

	// map otel data to columns
	for _, span := range spans {
		colTimestamp.Append(span.StartTime())
		colTraceId.Append(span.SpanContext().TraceID().String())
		colSpanId.Append(span.SpanContext().SpanID().String())
		colParentSpanId.Append(span.Parent().SpanID().String())
		colTraceState.Append(span.SpanContext().TraceState().String())
		colSpanName.Append(span.Name())
		colSpanKind.Append(span.SpanKind().String())
		colServiceName.Append(e.serviceName)

		resourceAttributes := make(map[string]string)
		for _, attr := range span.Resource().Attributes() {
			resourceAttributes[string(attr.Key)] = attr.Value.AsString()
		}
		colResourceAttributes.Append(resourceAttributes)

		colScopeName.Append(span.InstrumentationScope().Name)
		colScopeVersion.Append(span.InstrumentationScope().Version)

		spanAttributes := make(map[string]string)
		for _, attr := range span.Attributes() {
			spanAttributes[string(attr.Key)] = attr.Value.AsString()
		}
		colSpanAttributes.Append(spanAttributes)

		colDuration.Append(uint64(span.EndTime().Sub(span.StartTime())))
		colStatusCode.Append(span.Status().Code.String())
		colStatusMessage.Append(span.Status().Description)

		// map events
		var eventTimes []time.Time
		var eventNames []string
		var eventAttributes []map[string]string

		for _, event := range span.Events() {
			eventTimes = append(eventTimes, event.Time)
			eventNames = append(eventNames, event.Name)

			eventAttrs := make(map[string]string)
			for _, attr := range event.Attributes {
				eventAttrs[string(attr.Key)] = attr.Value.AsString()
			}

			eventAttributes = append(eventAttributes, eventAttrs)
		}

		colEventsTimestamp.Append(eventTimes)
		colEventsName.Append(eventNames)
		colEventsAttributes.Append(eventAttributes)

		// map links
		var linkTraceIds []string
		var linkSpanIds []string
		var linkTraceStates []string
		var linkAttributes []map[string]string

		for _, link := range span.Links() {
			linkTraceIds = append(linkTraceIds, link.SpanContext.TraceID().String())
			linkSpanIds = append(linkSpanIds, link.SpanContext.SpanID().String())
			linkTraceStates = append(linkTraceStates, link.SpanContext.TraceState().String())

			linkAttrs := make(map[string]string)
			for _, attr := range link.Attributes {
				linkAttrs[string(attr.Key)] = attr.Value.AsString()
			}

			linkAttributes = append(linkAttributes, linkAttrs)
		}

		colLinksTraceId.Append(linkTraceIds)
		colLinksSpanId.Append(linkSpanIds)
		colLinksTraceState.Append(linkTraceStates)
		colLinksAttributes.Append(linkAttributes)
	}

	// construct input block
	tracesInput := proto.Input{
		{Name: "Timestamp", Data: colTimestamp},
		{Name: "TraceId", Data: colTraceId},
		{Name: "SpanId", Data: colSpanId},
		{Name: "ParentSpanId", Data: colParentSpanId},
		{Name: "TraceState", Data: colTraceState},
		{Name: "SpanName", Data: colSpanName},
		{Name: "SpanKind", Data: colSpanKind},
		{Name: "ServiceName", Data: colServiceName},
		{Name: "ResourceAttributes", Data: colResourceAttributes},
		{Name: "ScopeName", Data: colScopeName},
		{Name: "ScopeVersion", Data: colScopeVersion},
		{Name: "SpanAttributes", Data: colSpanAttributes},
		{Name: "Duration", Data: colDuration},
		{Name: "StatusCode", Data: colStatusCode},
		{Name: "StatusMessage", Data: colStatusMessage},
		{Name: "Events.Timestamp", Data: colEventsTimestamp},
		{Name: "Events.Name", Data: colEventsName},
		{Name: "Events.Attributes", Data: colEventsAttributes},
		{Name: "Links.TraceId", Data: colLinksTraceId},
		{Name: "Links.SpanId", Data: colLinksSpanId},
		{Name: "Links.TraceState", Data: colLinksTraceState},
		{Name: "Links.Attributes", Data: colLinksAttributes},
	}

	// execute native protocol insert
	err := e.client.Do(ctx, ch.Query{
		Body:  `INSERT INTO otel_traces (Timestamp, TraceId, SpanId, ParentSpanId, TraceState, SpanName, SpanKind, ServiceName, ResourceAttributes, ScopeName, ScopeVersion, SpanAttributes, Duration, StatusCode, StatusMessage, "Events.Timestamp", "Events.Name", "Events.Attributes", "Links.TraceId", "Links.SpanId", "Links.TraceState", "Links.Attributes") VALUES`,
		Input: tracesInput,
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
