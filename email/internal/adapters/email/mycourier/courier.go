package mycourier

import (
	"context"
	"strconv"
	"time"

	"github.com/lordvidex/errs"
	"github.com/trycourier/courier-go/v2"
)

// Sender is a wrapper around courier client that implements the email.Sender interface
type Sender struct {
	c    *courier.Client
	vt   string // verificationTemplate
	rpt  string // resetPasswordTemplate
	nsst string // newSignInSessionTemplate
	exp  time.Duration
}

// New creates a new courier sender
func New(token string, opts ...func(*Sender)) (*Sender, error) {
	s := &Sender{}
	for _, opt := range opts {
		opt(s)
	}
	if token == "" {
		return nil, errs.B().Msg("CourierSender: Token is required").Err()
	}
	if s.vt == "" {
		return nil, errs.B().Msg("CourierSender: Verification template code is required").Err()
	}
	if s.rpt == "" {
		return nil, errs.B().Msg("CourierSender: Reset password template code is required").Err()
	}
	if s.nsst == "" {
		return nil, errs.B().Msg("CourierSender: New sign in session template code is required").Err()
	}
	if s.exp == 0 {
		return nil, errs.B().Msg("CourierSender: Expiration time is required").Err()
	}
	s.c = courier.CreateClient(token, nil)
	return s, nil
}

// WithExpiration sets the exp value in minutes
func WithExpiration(exp time.Duration) func(*Sender) {
	return func(s *Sender) {
		s.exp = exp
	}
}

// WithVerificationTemplate sets the verification template code
func WithVerificationTemplate(templateCode string) func(*Sender) {
	return func(s *Sender) {
		s.vt = templateCode
	}
}

// WithResetPasswordTemplate sets the reset password template code
func WithResetPasswordTemplate(templateCode string) func(*Sender) {
	return func(s *Sender) {
		s.rpt = templateCode
	}
}

// WithNewSignInSessionTemplate sets the new sign in session template code
func WithNewSignInSessionTemplate(templateCode string) func(*Sender) {
	return func(s *Sender) {
		s.nsst = templateCode
	}
}

// SendVerificationCode sends a verification code to the given email
func (c *Sender) SendVerificationCode(ctx context.Context, email string, name string, code string) error {
	requestID, err := c.c.SendMessage(ctx,
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to":       map[string]string{"email": email},
				"template": c.vt,
				"data": map[string]string{
					"name":       name,
					"code":       code,
					"expiration": strconv.Itoa(int(c.exp.Minutes())),
				},
			},
		},
	)
	if err != nil {
		return errs.B(err).Msgf("Failed to send verification code email, request ID: %s", requestID).Err()
	}
	return err
}

// SendResetPasswordToken sends a reset password token to the given email
// The token is used to reset the password and set a new one
func (c *Sender) SendResetPasswordToken(ctx context.Context, email string, name string, token string) error {
	requestID, err := c.c.SendMessage(ctx,
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to":       map[string]string{"email": email},
				"template": c.rpt,
				"data": map[string]string{
					"name":       name,
					"token":      token,
					"expiration": strconv.Itoa(int(c.exp.Minutes())),
				},
			},
		},
	)
	if err != nil {
		return errs.B(err).Msgf("Failed to send reset password token email, request ID: %s", requestID).Err()
	}
	return err
}

// SendNewSignInSession sends an email to notify user about a new login session on his account
func (c *Sender) SendNewSignInSession(ctx context.Context, email string, name string, clientIP string, userAgent string) error {
	requestID, err := c.c.SendMessage(ctx,
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to":       map[string]string{"email": email},
				"template": c.nsst,
				"data": map[string]string{
					"name":       name,
					"client_ip":  clientIP,
					"user_agent": userAgent,
				},
			},
		},
	)
	if err != nil {
		return errs.B(err).Msgf("Failed to send new signin session email, request ID: %s", requestID).Err()
	}
	return err
}

// Close closes the connection with the server
// Since the courier pkg doesn't have `close` function, this function returns nil
// This function is required to implement the `Sender` interface
func (c *Sender) Close() error {
	return nil
}
