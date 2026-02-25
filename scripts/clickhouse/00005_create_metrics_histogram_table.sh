#!/bin/sh
set -e

clickhouse-client --user oteldb_user --password "oteldb_pass" -n <<-EOSQL
  CREATE TABLE IF NOT EXISTS oteldb.otel_metrics_histogram
  (
    "TimeUnix" DateTime64(6) CODEC(Delta(8), LZ4),
    "ServiceName" LowCardinality(String) CODEC(LZ4),
    "MetricName" String CODEC(LZ4),
    "Count" UInt64 CODEC(Delta(8), LZ4),
    "Sum" Float64 CODEC(LZ4),
    "BucketCounts" Array(UInt64) CODEC(LZ4),
    "ExplicitBounds" Array(Float64) CODEC(LZ4),
    "Attributes" Map(LowCardinality(String), String) CODEC(LZ4),
    "ResourceAttributes" Map(LowCardinality(String), String) CODEC(LZ4),

    INDEX idx_res_attr_key mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_res_attr_value mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_attr_key mapKeys(Attributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_attr_value mapValues(Attributes) TYPE bloom_filter(0.01) GRANULARITY 1
  )
  ENGINE = MergeTree
  PARTITION BY toDate(TimeUnix)
  ORDER BY (ServiceName, MetricName, Attributes, TimeUnix)
  TTL toDate(TimeUnix) + toIntervalDay(30)
  SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;
EOSQL
