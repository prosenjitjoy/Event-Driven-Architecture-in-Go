package clickhouse

import (
	"context"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type MetricExporter struct {
	serviceName string
	client      *ch.Client
}

func NewMetricExporter(serviceName string, client *ch.Client) *MetricExporter {
	return &MetricExporter{
		serviceName: serviceName,
		client:      client,
	}
}

func (e *MetricExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	// initialize columns matching otel_metrics_gauge table schema
	colGaugeTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMicro)
	colGaugeServiceName := new(proto.ColStr).LowCardinality()
	colGaugeMetricName := new(proto.ColStr)
	colGaugeValue := new(proto.ColFloat64)
	colGaugeAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colGaugeResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))

	// initialize columns matching otel_metrics_histogram table schema
	colHistogramTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMicro)
	colHistogramServiceName := new(proto.ColStr).LowCardinality()
	colHistogramMetricName := new(proto.ColStr)
	colHistogramCount := new(proto.ColUInt64)
	colHistogramSum := new(proto.ColFloat64)
	colHistogramBucketCounts := new(proto.ColUInt64).Array()
	colHistogramExplicitBounds := new(proto.ColFloat64).Array()
	colHistogramAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colHistogramResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))

	// initialize columns matching otel_metrics_sum table schema
	colSumTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMicro)
	colSumServiceName := new(proto.ColStr).LowCardinality()
	colSumMetricName := new(proto.ColStr)
	colSumValue := new(proto.ColFloat64)
	colSumIsMonotonic := new(proto.ColBool)
	colSumAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colSumResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))

	// map otel data to columns
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Gauge[float64]:
				for _, dp := range data.DataPoints {
					colGaugeTimeUnix.Append(dp.Time)
					colGaugeServiceName.Append(e.serviceName)
					colGaugeMetricName.Append(m.Name)
					colGaugeValue.Append(dp.Value)

					gaugeAttributes := make(map[string]string)
					for _, attr := range dp.Attributes.ToSlice() {
						gaugeAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colGaugeAttributes.Append(gaugeAttributes)

					gaugeResourceAttributes := make(map[string]string)
					for _, attr := range rm.Resource.Attributes() {
						gaugeResourceAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colGaugeResourceAttributes.Append(gaugeResourceAttributes)
				}
			case metricdata.Histogram[float64]:
				for _, dp := range data.DataPoints {
					colHistogramTimeUnix.Append(dp.Time)
					colHistogramServiceName.Append(e.serviceName)
					colHistogramMetricName.Append(m.Name)
					colHistogramCount.Append(dp.Count)
					colHistogramSum.Append(dp.Sum)
					colHistogramBucketCounts.Append(dp.BucketCounts)
					colHistogramExplicitBounds.Append(dp.Bounds)

					histogramAttributes := make(map[string]string)
					for _, attr := range dp.Attributes.ToSlice() {
						histogramAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colHistogramAttributes.Append(histogramAttributes)

					histogramResourceAttributes := make(map[string]string)
					for _, attr := range rm.Resource.Attributes() {
						histogramResourceAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colHistogramResourceAttributes.Append(histogramResourceAttributes)
				}
			case metricdata.Sum[float64]:
				for _, dp := range data.DataPoints {
					colSumTimeUnix.Append(dp.Time)
					colSumServiceName.Append(e.serviceName)
					colSumMetricName.Append(m.Name)
					colSumValue.Append(dp.Value)
					colSumIsMonotonic.Append(data.IsMonotonic)

					sumAttributes := make(map[string]string)
					for _, attr := range dp.Attributes.ToSlice() {
						sumAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colSumAttributes.Append(sumAttributes)

					sumResourceAttributes := make(map[string]string)
					for _, attr := range rm.Resource.Attributes() {
						sumResourceAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colSumResourceAttributes.Append(sumResourceAttributes)
				}
			}
		}
	}

	// construct input block for metrics_gauge table
	inputGauge := proto.Input{
		{Name: "TimeUnix", Data: colGaugeTimeUnix},
		{Name: "ServiceName", Data: colGaugeServiceName},
		{Name: "MetricName", Data: colGaugeMetricName},
		{Name: "Value", Data: colGaugeValue},
		{Name: "Attributes", Data: colGaugeAttributes},
		{Name: "ResourceAttributes", Data: colGaugeResourceAttributes},
	}

	// construct input block for metrics_histogram table
	inputHistogram := proto.Input{
		{Name: "TimeUnix", Data: colHistogramTimeUnix},
		{Name: "ServiceName", Data: colHistogramServiceName},
		{Name: "MetricName", Data: colHistogramMetricName},
		{Name: "Count", Data: colHistogramCount},
		{Name: "Sum", Data: colHistogramSum},
		{Name: "BucketCounts", Data: colHistogramBucketCounts},
		{Name: "ExplicitBounds", Data: colHistogramExplicitBounds},
		{Name: "Attributes", Data: colHistogramAttributes},
		{Name: "ResourceAttributes", Data: colHistogramResourceAttributes},
	}

	// construct input block for metrics_sum table
	inputSum := proto.Input{
		{Name: "TimeUnix", Data: colSumTimeUnix},
		{Name: "ServiceName", Data: colSumServiceName},
		{Name: "MetricName", Data: colSumMetricName},
		{Name: "Value", Data: colSumValue},
		{Name: "IsMonotonic", Data: colSumIsMonotonic},
		{Name: "Attributes", Data: colSumAttributes},
		{Name: "ResourceAttributes", Data: colSumResourceAttributes},
	}

	// execute inserts only if columns are not empty
	if colGaugeTimeUnix.Rows() > 0 {
		err := e.client.Do(ctx, ch.Query{
			Body:  "INSERT INTO otel_metrics_gauge VALUES",
			Input: inputGauge,
		})
		if err != nil {
			return err
		}
	}
	if colHistogramTimeUnix.Rows() > 0 {
		err := e.client.Do(ctx, ch.Query{
			Body:  "INSERT INTO otel_metrics_histogram VALUES",
			Input: inputHistogram,
		})
		if err != nil {
			return err
		}
	}
	if colHistogramTimeUnix.Rows() > 0 {
		err := e.client.Do(ctx, ch.Query{
			Body:  "INSERT INTO otel_metrics_sum VALUES",
			Input: inputSum,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *MetricExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func (e *MetricExporter) Aggregation(kind metric.InstrumentKind) metric.Aggregation {
	return metric.DefaultAggregationSelector(kind)
}

func (e *MetricExporter) ForceFlush(ctx context.Context) error {
	return nil
}

func (e *MetricExporter) Shutdown(ctx context.Context) error {
	if e.client != nil {
		return e.client.Close()
	}

	return nil
}
