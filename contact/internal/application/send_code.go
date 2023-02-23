package application

import (
	"context"
	"time"

	"github.com/escalopa/gofly/contact/internal/core"
	"github.com/lordvidex/errs"
)

// SendCodeCommandParam is the parameter for SendCodeCommand
type SendCodeCommandParam struct {
	Email string
}

// SendCodeCommand is the command to send a verification code to a user
type SendCodeCommand interface {
	Execute(ctx context.Context, param SendCodeCommandParam) error
}

// SendCodeCommandImpl is the implementation of SendCodeCommand
type SendCodeCommandImpl struct {
	cr  CodeRepository
	cg  CodeGenerator
	es  EmailSender
	mti time.Duration // min time interval between sending codes
}

// NewSendCodeCommand creates a new SendCodeCommand
func NewSendCodeCommand(mti time.Duration, cr CodeRepository, cg CodeGenerator, es EmailSender) SendCodeCommand {
	return &SendCodeCommandImpl{
		cr:  cr,
		cg:  cg,
		es:  es,
		mti: mti,
	}
}

// Execute ctx context.Context, executes the command
func (c *SendCodeCommandImpl) Execute(ctx context.Context, param SendCodeCommandParam) error {
	// check if a message has been sent to the user in the last `c.mti`
	if vc, err := c.cr.Get(ctx, param.Email); err == nil {
		if time.Now().Sub(vc.SentAt) < c.mti {
			return errs.B().Msgf("please wait %d minute(s) before sending another code",
				int(time.Now().Add(c.mti).Sub(vc.SentAt).Minutes())).Err()
		}
	}
	// generate a new code
	code, err := c.cg.GenerateCode()
	if err != nil {
		return errs.B(err).Msg("could not generate code").Err()
	}
	// save the code to the database
	vc := core.VerificationCode{
		Code:   code,
		SentAt: time.Now(),
	}
	if err = c.cr.Save(ctx, param.Email, vc); err != nil {
		return errs.B(err).Msg("could not save code").Err()
	}
	// send the code to the user via email
	err = c.es.SendVerificationCode(ctx, param.Email, code)
	if err != nil {
		return errs.B(err).Msg("could not send email").Err()
	}
	return nil
}
