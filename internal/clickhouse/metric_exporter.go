package clickhouse

import (
	"context"
	"time"

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
	colGaugeResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colGaugeResourceSchemaUrl := new(proto.ColStr)
	colGaugeScopeName := new(proto.ColStr)
	colGaugeScopeVersion := new(proto.ColStr)
	colGaugeScopeAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colGaugeScopeDroppedAttrCount := new(proto.ColUInt32)
	colGaugeScopeSchemaUrl := new(proto.ColStr)
	colGaugeServiceName := new(proto.ColStr).LowCardinality()
	colGaugeMetricName := new(proto.ColStr)
	colGaugeMetricDescription := new(proto.ColStr)
	colGaugeMetricUnit := new(proto.ColStr)
	colGaugeAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colGaugeStartTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colGaugeTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colGaugeValue := new(proto.ColFloat64)
	colGaugeFlags := new(proto.ColUInt32)
	colGaugeExemplarsFilteredAttributes := proto.NewArray(proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr)))
	colGaugeExemplarsTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano).Array()
	colGaugeExemplarsValue := new(proto.ColFloat64).Array()
	colGaugeExemplarsSpanId := new(proto.ColStr).Array()
	colGaugeExemplarsTraceId := new(proto.ColStr).Array()

	// initialize columns matching otel_metrics_histogram table schema
	colHistogramResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colHistogramResourceSchemaUrl := new(proto.ColStr)
	colHistogramScopeName := new(proto.ColStr)
	colHistogramScopeVersion := new(proto.ColStr)
	colHistogramScopeAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colHistogramScopeDroppedAttrCount := new(proto.ColUInt32)
	colHistogramScopeSchemaUrl := new(proto.ColStr)
	colHistogramServiceName := new(proto.ColStr).LowCardinality()
	colHistogramMetricName := new(proto.ColStr)
	colHistogramMetricDescription := new(proto.ColStr)
	colHistogramMetricUnit := new(proto.ColStr)
	colHistogramAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colHistogramStartTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colHistogramTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colHistogramCount := new(proto.ColUInt64)
	colHistogramSum := new(proto.ColFloat64)
	colHistogramBucketCounts := new(proto.ColUInt64).Array()
	colHistogramExplicitBounds := new(proto.ColFloat64).Array()
	colHistogramExemplarsFilteredAttributes := proto.NewArray(proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr)))
	colHistogramExemplarsTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano).Array()
	colHistogramExemplarsValue := new(proto.ColFloat64).Array()
	colHistogramExemplarsSpanId := new(proto.ColStr).Array()
	colHistogramExemplarsTraceId := new(proto.ColStr).Array()
	colHistogramFlags := new(proto.ColUInt32)
	colHistogramMin := new(proto.ColFloat64)
	colHistogramMax := new(proto.ColFloat64)
	colHistogramAggregationTemporality := new(proto.ColInt32)

	// initialize columns matching otel_metrics_sum table schema
	colSumResourceAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colSumResourceSchemaUrl := new(proto.ColStr)
	colSumScopeName := new(proto.ColStr)
	colSumScopeVersion := new(proto.ColStr)
	colSumScopeAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colSumScopeDroppedAttrCount := new(proto.ColUInt32)
	colSumScopeSchemaUrl := new(proto.ColStr)
	colSumServiceName := new(proto.ColStr).LowCardinality()
	colSumMetricName := new(proto.ColStr)
	colSumMetricDescription := new(proto.ColStr)
	colSumMetricUnit := new(proto.ColStr)
	colSumAttributes := proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr))
	colSumStartTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colSumTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano)
	colSumValue := new(proto.ColFloat64)
	colSumFlags := new(proto.ColUInt32)
	colSumExemplarsFilteredAttributes := proto.NewArray(proto.NewMap(new(proto.ColStr).LowCardinality(), new(proto.ColStr)))
	colSumExemplarsTimeUnix := new(proto.ColDateTime64).WithPrecision(proto.PrecisionNano).Array()
	colSumExemplarsValue := new(proto.ColFloat64).Array()
	colSumExemplarsSpanId := new(proto.ColStr).Array()
	colSumExemplarsTraceId := new(proto.ColStr).Array()
	colSumAggregationTemporality := new(proto.ColInt32)
	colSumIsMonotonic := new(proto.ColBool)

	// map otel data to columns
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Gauge[float64]:
				for _, dp := range data.DataPoints {
					gaugeResourceAttributes := make(map[string]string)
					for _, attr := range rm.Resource.Attributes() {
						gaugeResourceAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colGaugeResourceAttributes.Append(gaugeResourceAttributes)

					colGaugeResourceSchemaUrl.Append(rm.Resource.SchemaURL())
					colGaugeScopeName.Append(sm.Scope.Name)
					colGaugeScopeVersion.Append(sm.Scope.Version)
					colGaugeScopeAttributes.Append(make(map[string]string))
					colGaugeScopeDroppedAttrCount.Append(0)
					colGaugeScopeSchemaUrl.Append(sm.Scope.SchemaURL)
					colGaugeServiceName.Append(e.serviceName)
					colGaugeMetricName.Append(m.Name)
					colGaugeMetricDescription.Append(m.Description)
					colGaugeMetricUnit.Append(m.Unit)

					gaugeAttributes := make(map[string]string)
					for _, attr := range dp.Attributes.ToSlice() {
						gaugeAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colGaugeAttributes.Append(gaugeAttributes)

					colGaugeStartTimeUnix.Append(dp.StartTime)
					colGaugeTimeUnix.Append(dp.Time)
					colGaugeValue.Append(dp.Value)
					colGaugeFlags.Append(0)

					exemplarFilteredAttritubes := make([]map[string]string, 0, len(dp.Exemplars))
					exemplarUnixTimes := make([]time.Time, 0, len(dp.Exemplars))
					exemplarValues := make([]float64, 0, len(dp.Exemplars))
					exemplarSpanIds := make([]string, 0, len(dp.Exemplars))
					exemplarTraceIds := make([]string, 0, len(dp.Exemplars))

					for _, exemplar := range dp.Exemplars {
						attrs := make(map[string]string)

						for _, kv := range exemplar.FilteredAttributes {
							attrs[string(kv.Key)] = kv.Value.AsString()
						}
						exemplarFilteredAttritubes = append(exemplarFilteredAttritubes, attrs)

						exemplarUnixTimes = append(exemplarUnixTimes, exemplar.Time)
						exemplarValues = append(exemplarValues, exemplar.Value)
						exemplarSpanIds = append(exemplarSpanIds, string(exemplar.SpanID))
						exemplarTraceIds = append(exemplarTraceIds, string(exemplar.TraceID))
					}

					colGaugeExemplarsFilteredAttributes.Append(exemplarFilteredAttritubes)
					colGaugeExemplarsTimeUnix.Append(exemplarUnixTimes)
					colGaugeExemplarsValue.Append(exemplarValues)
					colGaugeExemplarsSpanId.Append(exemplarSpanIds)
					colGaugeExemplarsTraceId.Append(exemplarTraceIds)
				}
			case metricdata.Histogram[float64]:
				for _, dp := range data.DataPoints {
					histogramResourceAttributes := make(map[string]string)
					for _, attr := range rm.Resource.Attributes() {
						histogramResourceAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colHistogramResourceAttributes.Append(histogramResourceAttributes)

					colHistogramResourceSchemaUrl.Append(rm.Resource.SchemaURL())
					colHistogramScopeName.Append(sm.Scope.Name)
					colHistogramScopeVersion.Append(sm.Scope.Version)
					colHistogramScopeAttributes.Append(make(map[string]string))
					colHistogramScopeDroppedAttrCount.Append(0)
					colHistogramScopeSchemaUrl.Append(sm.Scope.SchemaURL)
					colHistogramServiceName.Append(e.serviceName)
					colHistogramMetricName.Append(m.Name)
					colHistogramMetricDescription.Append(m.Description)
					colHistogramMetricUnit.Append(m.Unit)

					histogramAttributes := make(map[string]string)
					for _, attr := range dp.Attributes.ToSlice() {
						histogramAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colHistogramAttributes.Append(histogramAttributes)

					colHistogramStartTimeUnix.Append(dp.StartTime)
					colHistogramTimeUnix.Append(dp.Time)
					colHistogramCount.Append(dp.Count)
					colHistogramSum.Append(dp.Sum)
					colHistogramBucketCounts.Append(dp.BucketCounts)
					colHistogramExplicitBounds.Append(dp.Bounds)

					exemplarFilteredAttritubes := make([]map[string]string, 0, len(dp.Exemplars))
					exemplarUnixTimes := make([]time.Time, 0, len(dp.Exemplars))
					exemplarValues := make([]float64, 0, len(dp.Exemplars))
					exemplarSpanIds := make([]string, 0, len(dp.Exemplars))
					exemplarTraceIds := make([]string, 0, len(dp.Exemplars))

					for _, exemplar := range dp.Exemplars {
						attrs := make(map[string]string)

						for _, kv := range exemplar.FilteredAttributes {
							attrs[string(kv.Key)] = kv.Value.AsString()
						}
						exemplarFilteredAttritubes = append(exemplarFilteredAttritubes, attrs)

						exemplarUnixTimes = append(exemplarUnixTimes, exemplar.Time)
						exemplarValues = append(exemplarValues, exemplar.Value)
						exemplarSpanIds = append(exemplarSpanIds, string(exemplar.SpanID))
						exemplarTraceIds = append(exemplarTraceIds, string(exemplar.TraceID))
					}

					colHistogramExemplarsFilteredAttributes.Append(exemplarFilteredAttritubes)
					colHistogramExemplarsTimeUnix.Append(exemplarUnixTimes)
					colHistogramExemplarsValue.Append(exemplarValues)
					colHistogramExemplarsSpanId.Append(exemplarSpanIds)
					colHistogramExemplarsTraceId.Append(exemplarTraceIds)

					colHistogramFlags.Append(0)

					histogramMin, minDefined := dp.Min.Value()
					if minDefined {
						colHistogramMin.Append(histogramMin)
					} else {
						colHistogramMin.Append(0)
					}

					histogramMax, maxDefined := dp.Max.Value()
					if maxDefined {
						colHistogramMax.Append(histogramMax)
					} else {
						colHistogramMax.Append(0)
					}

					colHistogramAggregationTemporality.Append(int32(data.Temporality))
				}
			case metricdata.Sum[int64]:
				for _, dp := range data.DataPoints {
					sumResourceAttributes := make(map[string]string)
					for _, attr := range rm.Resource.Attributes() {
						sumResourceAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colSumResourceAttributes.Append(sumResourceAttributes)

					colSumResourceSchemaUrl.Append(rm.Resource.SchemaURL())
					colSumScopeName.Append(sm.Scope.Name)
					colSumScopeVersion.Append(sm.Scope.Version)
					colSumScopeAttributes.Append(make(map[string]string))
					colSumScopeDroppedAttrCount.Append(0)
					colSumScopeSchemaUrl.Append(sm.Scope.SchemaURL)
					colSumServiceName.Append(e.serviceName)
					colSumMetricName.Append(m.Name)
					colSumMetricDescription.Append(m.Description)
					colSumMetricUnit.Append(m.Unit)

					sumAttributes := make(map[string]string)
					for _, attr := range dp.Attributes.ToSlice() {
						sumAttributes[string(attr.Key)] = attr.Value.AsString()
					}
					colSumAttributes.Append(sumAttributes)

					colSumStartTimeUnix.Append(dp.StartTime)
					colSumTimeUnix.Append(dp.Time)
					colSumValue.Append(float64(dp.Value))
					colSumFlags.Append(0)

					exemplarFilteredAttritubes := make([]map[string]string, 0, len(dp.Exemplars))
					exemplarUnixTimes := make([]time.Time, 0, len(dp.Exemplars))
					exemplarValues := make([]float64, 0, len(dp.Exemplars))
					exemplarSpanIds := make([]string, 0, len(dp.Exemplars))
					exemplarTraceIds := make([]string, 0, len(dp.Exemplars))

					for _, exemplar := range dp.Exemplars {
						attrs := make(map[string]string)

						for _, kv := range exemplar.FilteredAttributes {
							attrs[string(kv.Key)] = kv.Value.AsString()
						}
						exemplarFilteredAttritubes = append(exemplarFilteredAttritubes, attrs)

						exemplarUnixTimes = append(exemplarUnixTimes, exemplar.Time)
						exemplarValues = append(exemplarValues, float64(exemplar.Value))
						exemplarSpanIds = append(exemplarSpanIds, string(exemplar.SpanID))
						exemplarTraceIds = append(exemplarTraceIds, string(exemplar.TraceID))
					}

					colSumExemplarsFilteredAttributes.Append(exemplarFilteredAttritubes)
					colSumExemplarsTimeUnix.Append(exemplarUnixTimes)
					colSumExemplarsValue.Append(exemplarValues)
					colSumExemplarsSpanId.Append(exemplarSpanIds)
					colSumExemplarsTraceId.Append(exemplarTraceIds)

					colSumAggregationTemporality.Append(int32(data.Temporality))
					colSumIsMonotonic.Append(data.IsMonotonic)
				}
			}
		}
	}

	// construct input block for metrics_gauge table
	gaugeInput := proto.Input{
		{Name: "ResourceAttributes", Data: colGaugeResourceAttributes},
		{Name: "ResourceSchemaUrl", Data: colGaugeResourceSchemaUrl},
		{Name: "ScopeName", Data: colGaugeScopeName},
		{Name: "ScopeVersion", Data: colGaugeScopeVersion},
		{Name: "ScopeAttributes", Data: colGaugeScopeAttributes},
		{Name: "ScopeDroppedAttrCount", Data: colGaugeScopeDroppedAttrCount},
		{Name: "ScopeSchemaUrl", Data: colGaugeScopeSchemaUrl},
		{Name: "ServiceName", Data: colGaugeServiceName},
		{Name: "MetricName", Data: colGaugeMetricName},
		{Name: "MetricDescription", Data: colGaugeMetricDescription},
		{Name: "MetricUnit", Data: colGaugeMetricUnit},
		{Name: "Attributes", Data: colGaugeAttributes},
		{Name: "StartTimeUnix", Data: colGaugeStartTimeUnix},
		{Name: "TimeUnix", Data: colGaugeTimeUnix},
		{Name: "Value", Data: colGaugeValue},
		{Name: "Flags", Data: colGaugeFlags},
		{Name: "Exemplars.FilteredAttributes", Data: colGaugeExemplarsFilteredAttributes},
		{Name: "Exemplars.TimeUnix", Data: colGaugeExemplarsTimeUnix},
		{Name: "Exemplars.Value", Data: colGaugeExemplarsValue},
		{Name: "Exemplars.SpanId", Data: colGaugeExemplarsSpanId},
		{Name: "Exemplars.TraceId", Data: colGaugeExemplarsTraceId},
	}

	// construct input block for metrics_histogram table
	histogramInput := proto.Input{
		{Name: "ResourceAttributes", Data: colHistogramResourceAttributes},
		{Name: "ResourceSchemaUrl", Data: colHistogramResourceSchemaUrl},
		{Name: "ScopeName", Data: colHistogramScopeName},
		{Name: "ScopeVersion", Data: colHistogramScopeVersion},
		{Name: "ScopeAttributes", Data: colHistogramScopeAttributes},
		{Name: "ScopeDroppedAttrCount", Data: colHistogramScopeDroppedAttrCount},
		{Name: "ScopeSchemaUrl", Data: colHistogramScopeSchemaUrl},
		{Name: "ServiceName", Data: colHistogramServiceName},
		{Name: "MetricName", Data: colHistogramMetricName},
		{Name: "MetricDescription", Data: colHistogramMetricDescription},
		{Name: "MetricUnit", Data: colHistogramMetricUnit},
		{Name: "Attributes", Data: colHistogramAttributes},
		{Name: "StartTimeUnix", Data: colHistogramStartTimeUnix},
		{Name: "TimeUnix", Data: colHistogramTimeUnix},
		{Name: "Count", Data: colHistogramCount},
		{Name: "Sum", Data: colHistogramSum},
		{Name: "BucketCounts", Data: colHistogramBucketCounts},
		{Name: "ExplicitBounds", Data: colHistogramExplicitBounds},
		{Name: "Exemplars.FilteredAttributes", Data: colHistogramExemplarsFilteredAttributes},
		{Name: "Exemplars.TimeUnix", Data: colHistogramExemplarsTimeUnix},
		{Name: "Exemplars.Value", Data: colHistogramExemplarsValue},
		{Name: "Exemplars.SpanId", Data: colHistogramExemplarsSpanId},
		{Name: "Exemplars.TraceId", Data: colHistogramExemplarsTraceId},
		{Name: "Flags", Data: colHistogramFlags},
		{Name: "Min", Data: colHistogramMin},
		{Name: "Max", Data: colHistogramMax},
		{Name: "AggregationTemporality", Data: colHistogramAggregationTemporality},
	}

	// construct input block for metrics_sum table
	sumInput := proto.Input{
		{Name: "ResourceAttributes", Data: colSumResourceAttributes},
		{Name: "ResourceSchemaUrl", Data: colSumResourceSchemaUrl},
		{Name: "ScopeName", Data: colSumScopeName},
		{Name: "ScopeVersion", Data: colSumScopeVersion},
		{Name: "ScopeAttributes", Data: colSumScopeAttributes},
		{Name: "ScopeDroppedAttrCount", Data: colSumScopeDroppedAttrCount},
		{Name: "ScopeSchemaUrl", Data: colSumScopeSchemaUrl},
		{Name: "ServiceName", Data: colSumServiceName},
		{Name: "MetricName", Data: colSumMetricName},
		{Name: "MetricDescription", Data: colSumMetricDescription},
		{Name: "MetricUnit", Data: colSumMetricUnit},
		{Name: "Attributes", Data: colSumAttributes},
		{Name: "StartTimeUnix", Data: colSumStartTimeUnix},
		{Name: "TimeUnix", Data: colSumTimeUnix},
		{Name: "Value", Data: colSumValue},
		{Name: "Flags", Data: colSumFlags},
		{Name: "Exemplars.FilteredAttributes", Data: colSumExemplarsFilteredAttributes},
		{Name: "Exemplars.TimeUnix", Data: colSumExemplarsTimeUnix},
		{Name: "Exemplars.Value", Data: colSumExemplarsValue},
		{Name: "Exemplars.SpanId", Data: colSumExemplarsSpanId},
		{Name: "Exemplars.TraceId", Data: colSumExemplarsTraceId},
		{Name: "AggregationTemporality", Data: colSumAggregationTemporality},
		{Name: "IsMonotonic", Data: colSumIsMonotonic},
	}

	// execute inserts only if columns are not empty
	if colGaugeTimeUnix.Rows() > 0 {
		err := e.client.Do(ctx, ch.Query{
			Body:  `INSERT INTO otel_metrics_gauge (ResourceAttributes, ResourceSchemaUrl, ScopeName, ScopeVersion, ScopeAttributes, ScopeDroppedAttrCount, ScopeSchemaUrl, ServiceName, MetricName, MetricDescription, MetricUnit, Attributes, StartTimeUnix, TimeUnix, Value, Flags, "Exemplars.FilteredAttributes", "Exemplars.TimeUnix", "Exemplars.Value", "Exemplars.SpanId", "Exemplars.TraceId") VALUES`,
			Input: gaugeInput,
		})
		if err != nil {
			return err
		}
	}
	if colHistogramTimeUnix.Rows() > 0 {
		err := e.client.Do(ctx, ch.Query{
			Body:  `INSERT INTO otel_metrics_histogram (ResourceAttributes, ResourceSchemaUrl, ScopeName, ScopeVersion, ScopeAttributes, ScopeDroppedAttrCount, ScopeSchemaUrl, ServiceName, MetricName, MetricDescription, MetricUnit, Attributes, StartTimeUnix, TimeUnix, Count, Sum, BucketCounts, ExplicitBounds, "Exemplars.FilteredAttributes", "Exemplars.TimeUnix", "Exemplars.Value", "Exemplars.SpanId", "Exemplars.TraceId", Flags, Min, Max, AggregationTemporality) VALUES`,
			Input: histogramInput,
		})
		if err != nil {
			return err
		}
	}
	if colSumTimeUnix.Rows() > 0 {
		err := e.client.Do(ctx, ch.Query{
			Body:  `INSERT INTO otel_metrics_sum (ResourceAttributes, ResourceSchemaUrl, ScopeName, ScopeVersion, ScopeAttributes, ScopeDroppedAttrCount, ScopeSchemaUrl, ServiceName, MetricName, MetricDescription, MetricUnit, Attributes, StartTimeUnix, TimeUnix, Value, Flags, "Exemplars.FilteredAttributes", "Exemplars.TimeUnix", "Exemplars.Value", "Exemplars.SpanId", "Exemplars.TraceId", AggregationTemporality, IsMonotonic) VALUES`,
			Input: sumInput,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *MetricExporter) Temporality(kind metric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
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
