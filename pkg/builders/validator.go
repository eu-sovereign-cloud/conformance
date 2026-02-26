package builders

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type requiredField struct {
	name  string
	value any
}

func field(name string, value any) requiredField {
	return requiredField{name: name, value: value}
}

func validateRequired(validate *validator.Validate, fields ...requiredField) error {
	for _, f := range fields {
		if err := validate.Var(f.value, "required"); err != nil {
			return fmt.Errorf("field %s is required", f.name)
		}
	}
	return nil
}

// TODO Find a better name for this function
func validateOneRequired(validate *validator.Validate, fields ...requiredField) error {
	errors := []error{}
	for _, f := range fields {
		if err := validate.Var(f.value, "required"); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == len(fields) {
		return errors[0]
	}
	return nil
}
