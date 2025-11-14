package secalib

import validator "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) ValidateRequired(field any) error {
	if err := v.validator.Var(field, "required"); err != nil {
		return err
	}
	return nil
}

func (v *Validator) ValidateRequireds(fields []any) error {
	for _, f := range fields {
		if err := v.validator.Var(f, "required"); err != nil {
			return err
		}
	}
	return nil
}
