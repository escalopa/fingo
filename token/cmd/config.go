package main

type config struct {
	GrpcPort string `mapstructure:"TOKEN_GRPC_PORT"`
	// Tracing
	TracingEnable            bool   `mapstructure:"TOKEN_TRACING_ENABLE"`
	TracingJaegerEnable      bool   `mapstructure:"TOKEN_TRACING_JAEGER_ENABLE"`
	TracingJaegerAgentUrl    string `mapstructure:"TOKEN_TRACING_JAEGER_AGENT_URL"`
	TracingJaegerServiceName string `mapstructure:"TOKEN_TRACING_JAEGER_SERVICE_NAME"`
	TracingJaegerEnvironment string `mapstructure:"TOKEN_TRACING_JAEGER_ENVIRONMENT"`
	// Tls
	GrpcTlsEnable   bool   `mapstructure:"TOKEN_GRPC_TLS_ENABLE"`
	GrpcTlsCertFile string `mapstructure:"TOKEN_GRPC_TLS_CERT_FILE"`
	GrpcTlsKeyFile  string `mapstructure:"TOKEN_GRPC_TLS_KEY_FILE"`
	// Storage
	RedisUrl string `mapstructure:"TOKEN_REDIS_URL"`
}

var cfg config
