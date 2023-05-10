package main

import "time"

type config struct {
	CodesExpiration                   time.Duration `mapstructure:"CONTACT_CODES_EXPIRATION"`
	SendCodeMinInterval               time.Duration `mapstructure:"CONTACT_SEND_CODE_MIN_INTERVAL"`
	SendResetPasswordTokenMinInterval time.Duration `mapstructure:"CONTACT_SEND_RESET_PASSWORD_TOKEN_MIN_INTERVAL"`
	// Tracing
	TracingEnable            bool   `mapstructure:"CONTACT_TRACING_ENABLE"`
	TracingJaegerEnable      bool   `mapstructure:"CONTACT_TRACING_JAEGER_ENABLE"`
	TracingJaegerAgentUrl    string `mapstructure:"CONTACT_TRACING_JAEGER_AGENT_URL"`
	TracingJaegerServiceName string `mapstructure:"CONTACT_TRACING_JAEGER_SERVICE_NAME"`
	TracingJaegerEnvironment string `mapstructure:"CONTACT_TRACING_JAEGER_ENVIRONMENT"`
	// RabbitMQ
	RabbitmqUrl                         string `mapstructure:"CONTACT_RABBITMQ_URL"`
	RabbitmqVerificationCodeQueueName   string `mapstructure:"CONTACT_RABBITMQ_VERIFICATION_CODE_QUEUE_NAME"`
	RabbitmqResetPasswordTokenQueueName string `mapstructure:"CONTACT_RABBITMQ_RESET_PASSWORD_TOKEN_QUEUE_NAME"`
	RabbitmqNewSigninSessionQueueName   string `mapstructure:"CONTACT_RABBITMQ_NEW_SIGNIN_SESSION_QUEUE_NAME"`
	// Couier
	CourierToken                      string `mapstructure:"CONTACT_COURIER_TOKEN"`
	CourierVerificationTemplateID     string `mapstructure:"CONTACT_COURIER_VERIFICATION_TEMPLATE_ID"`
	CourierResetPasswordTemplateID    string `mapstructure:"CONTACT_COURIER_RESET_PASSWORD_TEMPLATE_ID"`
	CourierNewSigninSessionTemplateID string `mapstructure:"CONTACT_COURIER_NEW_SIGNIN_SESSION_TEMPLATE_ID"`
}

var cfg config
