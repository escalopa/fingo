package tracer

import (
	"log"

	"github.com/lordvidex/errs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tr trace.Tracer
)

func init() {
	tr = newNoopTracer()
}

func SetTracer(t trace.Tracer) {
	tr = t
}

func Tracer() trace.Tracer {
	return tr
}

func newTracer(name string) trace.Tracer {
	// Create a new tracer provider
	tp := otel.Tracer(name)
	return tp
}

func newNoopTracer() trace.Tracer {
	return trace.NewNoopTracerProvider().Tracer("fingo")
}

func newJaegerExporter(url, service, env string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
			attribute.String("environment", env),
		)),
	)
	return tp, nil
}

// LoadTracer loads the tracer based on the environment variables
// tracingEnable: true if tracing is enabled, false otherwise
// tracingJaegerEnable: true if tracing is enable, false otherwise
// tracingJaegerAgentUrl: url of the jaeger agent
// tracingJaegerServiceName: name of the service
// tracingJaegerEnvironment: environment of the service
// If tracing is enabled, it will return a tracer otherwise, it will return a noop tracer
func LoadTracer(
	tracingEnable bool,
	tracingJaegerEnable bool,
	tracingJaegerAgentUrl string,
	tracingJaegerServiceName string,
	tracingJaegerEnvironment string,
) (trace.Tracer, error) {
	// check if tracing is enabled
	log.Println("starting server with tracing:", tracingEnable)
	if tracingEnable {
		// check if jaeger tracing is enabled
		log.Println("starting server with jaeger tracing:", tracingJaegerEnable)
		if tracingJaegerEnable {
			tp, err := newJaegerExporter(
				tracingJaegerAgentUrl,
				tracingJaegerServiceName,
				tracingJaegerEnvironment,
			)
			if err != nil {
				return nil, errs.B(err).Msg("failed to create jaeger exporter").Err()
			}
			otel.SetTracerProvider(tp)
		} else {
			// if jaeger is not enabled, use the default tracer
			otel.SetTracerProvider(tracesdk.NewTracerProvider())
		}
		return newTracer("fingo"), nil
	} else {
		return newNoopTracer(), nil
	}
}
