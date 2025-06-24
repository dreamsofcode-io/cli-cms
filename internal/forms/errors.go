package forms

import "errors"

var (
	// ErrUserCancelled is returned when the user cancels the form
	ErrUserCancelled = errors.New("user cancelled the operation")
	
	// ErrFormValidation is returned when form validation fails
	ErrFormValidation = errors.New("form validation failed")
)