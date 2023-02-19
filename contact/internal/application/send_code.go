package application

import (
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
	Execute(param SendCodeCommandParam) error
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

// Execute executes the command
func (c *SendCodeCommandImpl) Execute(param SendCodeCommandParam) error {
	// check if a message has been sent to the user in the last `c.mti`
	if vc, err := c.cr.Get(param.Email); err == nil {
		if time.Now().Sub(vc.SentAt) < c.mti {
			return errs.B().Msgf("please wait %d minute(s) before sending another code",
				time.Now().Sub(vc.SentAt).Minutes()).Err()
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
	if err = c.cr.Save(code, vc); err != nil {
		return errs.B(err).Msg("could not save code").Err()
	}
	// send the code to the user via email
	err = c.es.SendVerificationCode(param.Email, vc)
	if err != nil {
		return errs.B(err).Msg("could not send email").Err()
	}
	return nil
}
