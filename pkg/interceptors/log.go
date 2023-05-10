package interceptors

import (
	"context"
	"os"

	"github.com/escalopa/fingo/pkg/contextutils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02T15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

}

func LoggingUnaryInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		tID, _ := contextutils.GetTracerID(ctx)

		l := log.WithField("tracer-id", tID)

		// Add logger to context to have the same base field(tracer-id)
		ctx = contextutils.SetLogger(ctx, l)

		res, err := handler(ctx, req)
		if err != nil {
			l.Error(err.Error())
		}
		return res, err
	}
}
