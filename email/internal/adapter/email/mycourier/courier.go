package mycourier

import (
	"context"
	"strconv"
	"time"

	"github.com/lordvidex/errs"
	"github.com/trycourier/courier-go/v2"
)

type Sender struct {
	c   *courier.Client
	exp time.Duration
	vtc string //verificationTemplateCode
}

func New(token string, opts ...func(*Sender)) (*Sender, error) {
	s := &Sender{}
	for _, opt := range opts {
		opt(s)
	}
	if token == "" {
		return nil, errs.B().Msg("CourierSender: Token is required").Err()
	}
	if s.vtc == "" {
		return nil, errs.B().Msg("CourierSender: Verification template code is required").Err()
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

func WithVerificationTemplate(templateCode string) func(*Sender) {
	return func(s *Sender) {
		s.vtc = templateCode
	}
}

func (c *Sender) SendVerificationCode(ctx context.Context, email string, code string) error {
	requestID, err := c.c.SendMessage(
		ctx,
		courier.SendMessageRequestBody{
			Message: map[string]interface{}{
				"to": map[string]string{
					"email": email,
				},
				"template": c.vtc,
				"data": map[string]string{
					"code":       code,
					"expiration": strconv.Itoa(int(c.exp.Minutes())),
				},
			},
		},
	)
	if err != nil {
		return errs.B(err).Msgf("Failed to send verification code, request ID: %s", requestID).Err()
	}
	return err
}

// Close closes the connection with the server
// Since the courier pkg doesn't have `close` function, this function returns nil
func (c *Sender) Close() error {
	return nil
}
