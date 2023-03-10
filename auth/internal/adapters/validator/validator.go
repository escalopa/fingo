package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/lordvidex/errs"
)

type Validator struct {
	v *validator.Validate
}

type TestValidatorStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func NewValidator() *Validator {
	return &Validator{
		v: validator.New(),
	}
}

func (va *Validator) Validate(s any) error {
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
