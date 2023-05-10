package main

import "time"

type appConfig struct {
	GrpcPort string `mapstructure:"AUTH_GRPC_PORT"`
	// Tracing
	TracingEnable            bool   `mapstructure:"AUTH_TRACING_ENABLE"`
	TracingJaegerEnable      bool   `mapstructure:"AUTH_TRACING_JAEGER_ENABLE"`
	TracingJaegerAgentUrl    string `mapstructure:"AUTH_TRACING_JAEGER_AGENT_URL"`
	TracingJaegerServiceName string `mapstructure:"AUTH_TRACING_JAEGER_SERVICE_NAME"`
	TracingJaegerEnvironment string `mapstructure:"AUTH_TRACING_JAEGER_ENVIRONMENT"`
	// Tls
	GrpcTlsEnable   bool   `mapstructure:"AUTH_GRPC_TLS_ENABLE"`
	GrpcTlsCertFile string `mapstructure:"AUTH_GRPC_TLS_CERT_FILE"`
	GrpcTlsKeyFile  string `mapstructure:"AUTH_GRPC_TLS_KEY_FILE"`
	// Token Service
	TokenGrpcUrl             string `mapstructure:"TOKEN_GRPC_URL"`
	TokenGrpcTlsEnable       bool   `mapstructure:"TOKEN_GRPC_TLS_ENABLE"`
	TokenGrpcTlsUserCertFile string `mapstructure:"AUTH_TOKEN_GRPC_TLS_USER_CERT_FILE"`
	// Authentication
	TokenSecret          string        `mapstructure:"AUTH_TOKEN_SECRET"`
	AccessTokenDuration  time.Duration `mapstructure:"AUTH_ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"AUTH_REFRESH_TOKEN_DURATION"`
	UserSessionDuration  time.Duration `mapstructure:"AUTH_USER_SESSION_DURATION"`
	// Storage
	RedisUrl              string `mapstructure:"AUTH_CACHE_URL"`
	DatabaseMigrationPath string `mapstructure:"AUTH_DATABASE_MIGRATION_PATH"`
	DatabaseUrl           string `mapstructure:"AUTH_DATABASE_URL"`
	// Rabbitmq
	RabbitmqUrl                       string `mapstructure:"AUTH_RABBITMQ_URL"`
	RabbitmqNewSigninSessionQueueName string `mapstructure:"AUTH_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME"`
}

var cfg appConfig
