#!/bin/sh
set -e

clickhouse-client --user oteldb_user --password "oteldb_pass" -n <<-EOSQL
  CREATE TABLE IF NOT EXISTS oteldb.otel_logs
  (
    "Timestamp" DateTime64(6) CODEC(Delta(8), LZ4),
    "TraceId" String CODEC(LZ4),
    "SpanId" String CODEC(LZ4),
    "SeverityText" LowCardinality(String) CODEC(LZ4),
    "ServiceName" LowCardinality(String) CODEC(LZ4),
    "Body" String CODEC(LZ4),
    "ResourceAttributes" Map(LowCardinality(String), String) CODEC(LZ4),
    "LogAttributes" Map(LowCardinality(String), String) CODEC(LZ4),

    INDEX idx_trace_id TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_res_attr_key mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_res_attr_value mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_log_attr_key mapKeys(LogAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_log_attr_value mapValues(LogAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_lower_body lower(Body) TYPE tokenbf_v1(32768, 3, 0) GRANULARITY 8
  )
  ENGINE = MergeTree
  PARTITION BY toDate(Timestamp)
  ORDER BY (ServiceName, toDateTime(Timestamp), TraceId)
  TTL toDate(Timestamp) + toIntervalDay(30)
  SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;
EOSQL
