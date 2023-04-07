package pkgvalidator

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/lordvidex/errs"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

func (va *Validator) Validate(_ context.Context, s any) error {
	err := va.v.Struct(s)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			es := make([]string, len(errors))
			for i, e := range errors {
				es[i] = e.Error()
			}
			return errs.B().Code(errs.InvalidArgument).Msg(es...).Err()
		} else {
			return errs.B().Code(errs.Internal).Msg(err.Error()).Err()
		}
	}
	return nil
}
