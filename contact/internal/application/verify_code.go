package application

import "github.com/lordvidex/errs"

type VerifyCodeCommandParam struct {
	Email string
	Code  string
}

// VerifyCodeCommand is the interface for the VerifyCodeCommand
// It is used to verify a code sent to an email
type VerifyCodeCommand interface {
	Execute(param VerifyCodeCommandParam) error
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
func (c *VerifyCodeCommandImpl) Execute(param VerifyCodeCommandParam) error {
	// verify code that it does not contain any special characters & is not empty
	if !c.cg.VerifyCode(param.Code) {
		return errs.B().Msg("invalid code").Err()
	}
	// get the email associated with the code
	vc, err := c.cr.Get(param.Email)
	if err != nil {
		// if the code does not exist, return a not found error
		errr, ok := err.(*errs.Error)
		if ok && errr.Code == errs.NotFound {
			return errs.B(errr).Msg("code has expired").Err()
		}
		// otherwise, return the error
		return err
	}
	// if the email does not match the email in the param, return an error
	if vc.Code != param.Code {
		return errs.B().Msgf("code & email mismatch, expected %s, got %s", param.Code, vc.Code).Err()
	}
	return nil
}
