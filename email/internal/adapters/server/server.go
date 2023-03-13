package server

import (
	"context"

	"github.com/escalopa/fingo/email/internal/core"

	"github.com/escalopa/fingo/email/internal/application"
)

type Server struct {
	uc      *application.UseCases
	mqc     application.MessageConsumer
	errChan chan error
}

func NewServer(uc *application.UseCases, mqc application.MessageConsumer) *Server {
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
	err := s.mqc.HandleSendVerificationsCode(func(ctx context.Context, params core.SendVerificationCodeMessage) error {
		return s.uc.SendVerificationCode.Execute(ctx, application.SendVerificationCodeCommandParam{
			Name:  params.Name,
			Email: params.Email,
			Code:  params.Code,
		})
	})
	return err
}

func (s *Server) handleSendResetPasswordToken() error {
	err := s.mqc.HandleSendResetPasswordToken(func(ctx context.Context, params core.SendResetPasswordTokenMessage) error {
		return s.uc.SendResetPasswordToken.Execute(ctx, application.SendResetPasswordTokenCommandParam{
			Name:  params.Name,
			Email: params.Email,
			Token: params.Token,
		})
	})
	return err
}

func (s *Server) handleSendNewSignInSessionCode() error {
	err := s.mqc.HandleSendNewSignInSession(func(ctx context.Context, params core.SendNewSignInSessionMessage) error {
		return s.uc.SendNewSignInSession.Execute(ctx, application.SendNewSingInSessionCommandParam{
			Name:      params.Name,
			Email:     params.Email,
			ClientIP:  params.ClientIP,
			UserAgent: params.UserAgent,
		})
	})
	return err
}
