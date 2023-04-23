package oteltracer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestSetTracer(t *testing.T) {

	tests := []struct {
		name string
		want trace.Tracer
		set  trace.Tracer
	}{
		{
			name: "noop tracer",
			want: trace.NewNoopTracerProvider().Tracer("fingo"),
			set:  trace.NewNoopTracerProvider().Tracer("fingo"),
		},
		{
			name: "create tracer",
			want: otel.Tracer("fingo"),
			set:  otel.Tracer("fingo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTracer(tt.set)
			require.Equal(t, tt.want, Tracer())
		})
	}
}
