package server

import (
	"context"
	"log"

	"github.com/escalopa/fingo/contact/internal/core"

	"github.com/escalopa/fingo/contact/internal/application"
)

type Server struct {
	uc   *application.UseCases
	cons application.MessageConsumer

	errChan  chan error
	exitChan chan struct{}
}

func NewServer(uc *application.UseCases, cons application.MessageConsumer) *Server {
	return &Server{
		uc:   uc,
		cons: cons,

		errChan:  make(chan error, 1),
		exitChan: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	defer close(s.exitChan)
	defer close(s.errChan)

	// Notify err channel on error when starting handlers
	go func() {
		select {
		case s.errChan <- s.handleSendEmailVerificationCode():
		case s.errChan <- s.handleSendResetPasswordToken():
		case s.errChan <- s.handleSendNewSignInSessionCode():
		}
	}()

	// Wait for handlers throw error or exit signal
	select {
	case err := <-s.errChan:
		return err
	case <-s.exitChan:
		return nil
	}
}

// GracefulStop stops the consumer and returns
func (s *Server) GracefulStop() {
	err := s.cons.Close()
	if err != nil {
		log.Println("failed to stop consumer ", err)
	}
}

// Stop sends a signal to the server to stop
func (s *Server) Stop() {
	s.exitChan <- struct{}{}
}

func (s *Server) handleSendEmailVerificationCode() error {
	err := s.cons.HandleSendVerificationsCode(func(ctx context.Context, params core.SendVerificationCodeMessage) error {
		return s.uc.SendVerificationCode.Execute(ctx, application.SendVerificationCodeCommandParam{
			Name:  params.Name,
			Email: params.Email,
			Code:  params.Code,
		})
	})
	return err
}

func (s *Server) handleSendResetPasswordToken() error {
	err := s.cons.HandleSendResetPasswordToken(func(ctx context.Context, params core.SendResetPasswordTokenMessage) error {
		return s.uc.SendResetPasswordToken.Execute(ctx, application.SendResetPasswordTokenCommandParam{
			Name:  params.Name,
			Email: params.Email,
			Token: params.Token,
		})
	})
	return err
}

func (s *Server) handleSendNewSignInSessionCode() error {
	err := s.cons.HandleSendNewSignInSession(func(ctx context.Context, params core.SendNewSignInSessionMessage) error {
		return s.uc.SendNewSignInSession.Execute(ctx, application.SendNewSingInSessionCommandParam{
			Name:      params.Name,
			Email:     params.Email,
			ClientIP:  params.ClientIP,
			UserAgent: params.UserAgent,
		})
	})
	return err
}
