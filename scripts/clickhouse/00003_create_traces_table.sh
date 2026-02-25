#!/bin/sh
set -e

clickhouse-client --user oteldb_user --password "oteldb_pass" -n <<-EOSQL
  CREATE TABLE IF NOT EXISTS oteldb.otel_traces
  (
    "Timestamp" DateTime64(6) CODEC(Delta(8), LZ4),
    "TraceId" String CODEC(LZ4),
    "SpanId" String CODEC(LZ4),
    "ParentSpanId" String CODEC(LZ4),
    "SpanName" LowCardinality(String) CODEC(LZ4),
    "SpanKind" LowCardinality(String) CODEC(LZ4),
    "ServiceName" LowCardinality(String) CODEC(LZ4),
    "ResourceAttributes" Map(LowCardinality(String), String) CODEC(LZ4),
    "SpanAttributes" Map(LowCardinality(String), String) CODEC(LZ4),
    "Duration" UInt64 CODEC(LZ4),
    "StatusCode" LowCardinality(String) CODEC(LZ4),
    "StatusMessage" String CODEC(LZ4),
      
    INDEX idx_trace_id TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_res_attr_key mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_res_attr_value mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_span_attr_key mapKeys(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_span_attr_value mapValues(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_duration Duration TYPE minmax GRANULARITY 1
  )
  ENGINE = MergeTree
  PARTITION BY toDate(Timestamp)
  ORDER BY (ServiceName, SpanName, toDateTime(Timestamp))
  TTL toDate(Timestamp) + toIntervalDay(30)
  SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;
EOSQL
