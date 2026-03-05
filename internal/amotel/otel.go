package amotel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	propagator propagation.TextMapPropagator
	tracer     trace.Tracer
	meter      metric.Meter
)

func init() {
	propagator = otel.GetTextMapPropagator()
	tracer = otel.Tracer("internal/amotel")
	meter = otel.Meter("internal/amotel")
}
