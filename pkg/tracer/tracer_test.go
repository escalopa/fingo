package tracer

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

func TestNewJaegerExporter(t *testing.T) {
	type arg struct {
		tracingJaegerAgentUrl    string
		tracingJaegerServiceName string
		tracingJaegerEnvironment string
	}
	tests := []struct {
		name    string
		arg     arg
		wantErr bool
	}{
		{
			name: "success",
			arg: arg{
				tracingJaegerAgentUrl:    testJagerUrl,
				tracingJaegerServiceName: "fingo",
				tracingJaegerEnvironment: "test",
			},
			wantErr: false,
		},
		{
			// TODO: Know why this test passes despite url is emoty
			name: "empty url",
			arg: arg{
				tracingJaegerAgentUrl:    "",
				tracingJaegerServiceName: "fingo",
				tracingJaegerEnvironment: "test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newJaegerExporter(tt.arg.tracingJaegerAgentUrl, tt.arg.tracingJaegerServiceName, tt.arg.tracingJaegerEnvironment)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestLoadTracer(t *testing.T) {
	type arg struct {
		tracingEnable            bool
		tracingJaegerEnable      bool
		tracingJaegerAgentUrl    string
		tracingJaegerServiceName string
		tracingJaegerEnvironment string
	}

	tests := []struct {
		name    string
		arg     arg
		wantErr bool
	}{
		{
			name: "success enable tracing",
			arg: arg{
				tracingEnable:            true,
				tracingJaegerEnable:      true,
				tracingJaegerAgentUrl:    testJagerUrl,
				tracingJaegerServiceName: "fingo",
				tracingJaegerEnvironment: "test",
			},
			wantErr: false,
		},
		{
			name: "success disable tracing",
			arg: arg{
				tracingEnable:            false,
				tracingJaegerEnable:      false,
				tracingJaegerAgentUrl:    "",
				tracingJaegerServiceName: "",
				tracingJaegerEnvironment: "",
			},
			wantErr: false,
		},
		{
			name: "success disable jaeger",
			arg: arg{
				tracingEnable:            true,
				tracingJaegerEnable:      false,
				tracingJaegerAgentUrl:    "",
				tracingJaegerServiceName: "",
				tracingJaegerEnvironment: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadTracer(tt.arg.tracingEnable, tt.arg.tracingJaegerEnable, tt.arg.tracingJaegerAgentUrl, tt.arg.tracingJaegerServiceName, tt.arg.tracingJaegerEnvironment)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}
