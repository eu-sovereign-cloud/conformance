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

func (v *Validator) ValidateRequired(fields ...any) error {
	for _, f := range fields {
		if err := v.validator.Var(f, "required"); err != nil {
			return err
		}
	}
	return nil
}

// TODO Find a better name for this function
func (v *Validator) ValidateOneRequired(fields ...any) error {
	errors := []error{}
	for _, f := range fields {
		if err := v.validator.Var(f, "required"); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == len(fields) {
		return errors[0]
	}
	return nil
}
