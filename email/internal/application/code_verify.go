package application

import (
	"context"
	"github.com/lordvidex/errs"
)

type VerifyCodeCommandParam struct {
	Email string
	Code  string
}

// VerifyCodeCommand is the interface for the VerifyCodeCommand
// It is used to verify a code sent to an email
type VerifyCodeCommand interface {
	Execute(ctx context.Context, param VerifyCodeCommandParam) error
}

// VerifyCodeCommandImpl is the implementation of VerifyCodeCommand
type VerifyCodeCommandImpl struct {
	cr CodeRepository
	cg CodeGenerator
}

// NewVerifyCodeCommand creates a new VerifyCodeCommand
func NewVerifyCodeCommand(cr CodeRepository, cg CodeGenerator) VerifyCodeCommand {
	return &VerifyCodeCommandImpl{
		cr: cr,
		cg: cg,
	}
}

// Execute executes the command & verifies the code
func (c *VerifyCodeCommandImpl) Execute(ctx context.Context, param VerifyCodeCommandParam) error {
	// verify code that it does not contain any special characters & is not empty
	if !c.cg.VerifyCode(param.Code) {
		return errs.B().Msg("invalid code").Err()
	}
	// get the email associated with the code
	vc, err := c.cr.Get(ctx, param.Email)
	if err != nil {
		// if the code does not exist, return a not found error
		er, ok := err.(*errs.Error)
		if ok && er.Code == errs.NotFound {
			return errs.B(er).Msg("code has expired").Err()
		}
		// otherwise, return the error
		return err
	}
	// if the email does not match the email in the param, return an error
	if vc.Code != param.Code {
		return errs.B().Msg("given code doesn't match the one stored in cache").Err()
	}
	return nil
}
