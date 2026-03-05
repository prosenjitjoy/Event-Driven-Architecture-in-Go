package clickhouse

import (
	"context"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	logs "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
)

type LogExporter struct {
	serviceName string
	client      *ch.Client
}

func NewLogExpoerter(serviceName string, client *ch.Client) *LogExporter {
	return &LogExporter{
		serviceName: serviceName,
		client:      client,
	}
}

func (e *LogExporter) Export(ctx context.Context, records []log.Record) error {
	// initialize columns matching otel_logs table schema
	colTimestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colTraceId := new(proto.ColStr)
	colSpanId := new(proto.ColStr)
	colTraceFlags := new(proto.ColUInt8)
	colSeverityText := new(proto.ColStr).LowCardinality()
	colSeverityNumber := new(proto.ColUInt8)
	colServiceName := new(proto.ColStr).LowCardinality()
	colBody := new(proto.ColStr)
	colResourceSchemaUrl := new(proto.ColStr).LowCardinality()
	colResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colScopeSchemaUrl := new(proto.ColStr).LowCardinality()
	colScopeName := new(proto.ColStr)
	colScopeVersion := new(proto.ColStr).LowCardinality()
	colScopeAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colLogAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))

	var totalLogAttrCount uint64 = 0

	// map otel data to columns
	for _, record := range records {
		colTimestamp.Append(record.Timestamp())
		colTraceId.Append(record.TraceID().String())
		colSpanId.Append(record.SpanID().String())
		colTraceFlags.Append(uint8(record.TraceFlags()))
		colSeverityText.Append(record.SeverityText())
		colSeverityNumber.Append(uint8(record.Severity()))
		colServiceName.Append(e.serviceName)
		colBody.Append(record.Body().AsString())
		colResourceSchemaUrl.Append(record.Resource().SchemaURL())

		resourceAttributes := make(map[string]string)
		for _, attr := range record.Resource().Attributes() {
			resourceAttributes[string(attr.Key)] = attr.Value.AsString()
		}
		colResourceAttributes.Append(resourceAttributes)

		colScopeSchemaUrl.Append(record.InstrumentationScope().SchemaURL)
		colScopeName.Append(record.InstrumentationScope().Name)
		colScopeVersion.Append(record.InstrumentationScope().Version)
		colScopeAttributes.Append(make(map[string]string))

		// get reference to the internal columns of the map to be used inside callback
		logKeys := colLogAttributes.Keys.(*proto.ColLowCardinality[string])
		logValues := colLogAttributes.Values.(*proto.ColStr)

		var logAttrCount uint64 = 0

		record.WalkAttributes(func(kv logs.KeyValue) bool {
			logKeys.Append(kv.Key)
			logValues.Append(kv.Value.AsString())
			logAttrCount++
			return true
		})

		// update the cumulative offset
		totalLogAttrCount += logAttrCount
		colLogAttributes.Offsets.Append(totalLogAttrCount)
	}

	// construct input block
	logsInput := proto.Input{
		{Name: "Timestamp", Data: colTimestamp},
		{Name: "TraceId", Data: colTraceId},
		{Name: "SpanId", Data: colSpanId},
		{Name: "TraceFlags", Data: colTraceFlags},
		{Name: "SeverityText", Data: colSeverityText},
		{Name: "SeverityNumber", Data: colSeverityNumber},
		{Name: "ServiceName", Data: colServiceName},
		{Name: "Body", Data: colBody},
		{Name: "ResourceSchemaUrl", Data: colResourceSchemaUrl},
		{Name: "ResourceAttributes", Data: colResourceAttributes},
		{Name: "ScopeSchemaUrl", Data: colScopeSchemaUrl},
		{Name: "ScopeName", Data: colScopeName},
		{Name: "ScopeVersion", Data: colScopeVersion},
		{Name: "ScopeAttributes", Data: colScopeAttributes},
		{Name: "LogAttributes", Data: colLogAttributes},
	}

	// execute native protocol insert
	err := e.client.Do(ctx, ch.Query{
		Body:  `INSERT INTO otel_logs (Timestamp, TraceId, SpanId, TraceFlags, SeverityText, SeverityNumber, ServiceName, Body, ResourceSchemaUrl, ResourceAttributes, ScopeSchemaUrl, ScopeName, ScopeVersion, ScopeAttributes, LogAttributes) VALUES`,
		Input: logsInput,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *LogExporter) ForceFlush(ctx context.Context) error {
	return nil
}

func (e *LogExporter) Shutdown(ctx context.Context) error {
	if e.client != nil {
		return e.client.Close()
	}

	return nil
}
