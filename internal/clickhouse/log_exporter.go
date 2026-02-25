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
	colTimestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMicro)
	colTraceId := new(proto.ColStr)
	colSpanId := new(proto.ColStr)
	colSeverityText := new(proto.ColStr).LowCardinality()
	colServiceName := new(proto.ColStr).LowCardinality()
	colBody := new(proto.ColStr)
	colResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colLogAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))

	var totalLogAttrCount uint64 = 0

	// map otel data to columns
	for _, record := range records {
		colTimestamp.Append(record.Timestamp())
		colTraceId.Append(record.TraceID().String())
		colSpanId.Append(record.SpanID().String())
		colSeverityText.Append(record.SeverityText())
		colServiceName.Append(e.serviceName)
		colBody.Append(record.Body().AsString())

		resourceAttributes := make(map[string]string)
		for _, attr := range record.Resource().Attributes() {
			resourceAttributes[string(attr.Key)] = attr.Value.AsString()
		}
		colResourceAttributes.Append(resourceAttributes)

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
	input := proto.Input{
		{Name: "Timestamp", Data: colTimestamp},
		{Name: "TraceId", Data: colTraceId},
		{Name: "SpanId", Data: colSpanId},
		{Name: "SeverityText", Data: colSeverityText},
		{Name: "ServiceName", Data: colServiceName},
		{Name: "Body", Data: colBody},
		{Name: "ResourceAttributes", Data: colResourceAttributes},
		{Name: "LogAttributes", Data: colLogAttributes},
	}

	// execute native protocol insert
	err := e.client.Do(ctx, ch.Query{
		Body:  "INSERT INTO otel_logs VALUES",
		Input: input,
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
