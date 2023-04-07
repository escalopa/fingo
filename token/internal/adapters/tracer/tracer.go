package oteltracer

import "go.opentelemetry.io/otel/trace"

var (
	tr trace.Tracer
)

func init() {
	tr = trace.NewNoopTracerProvider().Tracer("fingo")
}

func SetTracer(t trace.Tracer) {
	tr = t
}

func Tracer() trace.Tracer {
	return tr
}
