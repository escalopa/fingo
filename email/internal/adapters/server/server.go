package server

import (
	"context"

	"github.com/escalopa/fingo/email/internal/application"
)

type Server struct {
	uc      *application.UseCases
	mqc     application.MessageQueueConsumer
	errChan chan error
}

func NewServer(uc *application.UseCases, mqc application.MessageQueueConsumer) *Server {
	return &Server{
		uc:  uc,
		mqc: mqc,
	}
}

func (s *Server) Start() error {
	go func() { s.errChan <- s.handleSendEmailVerificationCode() }()
	go func() { s.errChan <- s.handleSendResetPasswordToken() }()
	go func() { s.errChan <- s.handleSendNewSignInSessionCode() }()
	return <-s.errChan
}

func (s *Server) Stop() {
	s.errChan <- nil
	close(s.errChan)
}

func (s *Server) handleSendEmailVerificationCode() error {
	err := s.mqc.HandleSendVerificationsCode(func(ctx context.Context, email string, code string) error {
		return s.uc.SendVerificationCode.Execute(ctx, application.SendVerificationCodeCommandParam{
			Email: email,
			Code:  code,
		})
	})
	return err
}

func (s *Server) handleSendResetPasswordToken() error {
	err := s.mqc.HandleSendResetPasswordToken(func(ctx context.Context, email string, token string) error {
		return s.uc.SendResetPasswordToken.Execute(ctx, application.SendResetPasswordTokenCommandParam{
			Email: email,
			Token: token,
		})
	})
	return err
}

func (s *Server) handleSendNewSignInSessionCode() error {
	err := s.mqc.HandleSendNewSignInSession(func(ctx context.Context, email string, clientIP string, userAgent string) error {
		return s.uc.SendNewSignInSession.Execute(ctx, application.SendNewSingInSessionCommandParam{
			Email:     email,
			ClientIP:  clientIP,
			UserAgent: userAgent,
		})
	})
	return err
}
