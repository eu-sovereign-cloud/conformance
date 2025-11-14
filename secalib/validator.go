package secalib

import validator "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func newValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) ValidateRequireds(fields []any) error {
	for _, f := range fields {
		if err := v.validator.Var(f, "required"); err != nil {
			return err
		}
	}
	return nil
}
