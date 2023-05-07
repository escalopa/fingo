package main

import "time"

type config struct {
	GrpcPort     string `mapstructure:"WALLET_GRPC_PORT"`
	TokenGrpcUrl string `mapstructure:"TOKEN_GRPC_URL"`
	// Tracing
	TracingEnable            bool   `mapstructure:"WALLET_TRACING_ENABLE"`
	TracingJaegerEnable      bool   `mapstructure:"WALLET_TRACING_JAEGER_ENABLE"`
	TracingJaegerAgentUrl    string `mapstructure:"WALLET_TRACING_JAEGER_AGENT_URL"`
	TracingJaegerServiceName string `mapstructure:"WALLET_TRACING_JAEGER_SERVICE_NAME"`
	TracingJaegerEnvironment string `mapstructure:"WALLET_TRACING_JAEGER_ENVIRONMENT"`
	// Tls
	GrpcTlsEnable   bool   `mapstructure:"WALLET_GRPC_TLS_ENABLE"`
	GrpcTlsCertFile string `mapstructure:"WALLET_GRPC_TLS_CERT_FILE"`
	GrpcTlsKeyFile  string `mapstructure:"WALLET_GRPC_TLS_KEY_FILE"`
	// Token server
	TokenGrpcTlsEnable       bool   `mapstructure:"TOKEN_GRPC_TLS_ENABLE"`
	TokenGrpcTlsUserCertFile string `mapstructure:"WALLET_TOKEN_GRPC_TLS_USER_CERT_FILE"`
	// Storage
	DatabaseUrl           string `mapstructure:"WALLET_DATABASE_URL"`
	DatabaseMigrationPath string `mapstructure:"WALLET_DATABASE_MIGRATION_PATH"`
	// Card generation
	CardNumberLength int `mapstructure:"WALLET_CARD_NUMBER_LENGTH"`
	// Locker
	LockerCleanupDuration time.Duration `mapstructure:"WALLET_LOCKER_CLEANUP_DURATION"`
}

var cfg config
