package mcclintock

import (
	"strings"
)

// ValidationFunc is a type alias for a validator that takes a string and
// returns true if it passes validation.
type ValidationFunc func(s string) bool

// ValidateYesNo is a validation function that ensures that a yes/no question
// receives a valid response.  Use IsYes or IsNo to test for a particular
// yes or no response.
func ValidateYesNo(s string) bool {
	options := []string{"yes", "y", "no", "n"}
	return validateOptions(strings.ToLower(s), options)
}

// IsYes returns true if the response is some form of yes or y
func IsYes(s string) bool {
	options := []string{"yes", "y"}
	return validateOptions(strings.ToLower(s), options)
}

// IsNo returns true if the response is some form of no or n
func IsNo(s string) bool {
	options := []string{"no", "n"}
	return validateOptions(strings.ToLower(s), options)
}

// NewChoiceValidator builds a new validator for an enumerated options
// list.  This is used by the Choice function to validate a users response
// from a number of options.
//func NewChoiceValidator(options map[string]string) ValidationFunc {
//
//}

// NewOptionValidator builds a new validator from a []string of options
func NewOptionValidator(options []string) ValidationFunc {
	return func(s string) bool {
		return validateOptions(strings.ToLower(s), options)
	}
}

func ValidateWithFunc(s string, f ValidationFunc) bool {
	if f(s) {
		return true
	}
	return false
}

// validateOptions returns true if the given string appears in the list of valid
// options.
func validateOptions(s string, options []string) bool {
	l := strings.ToLower(s)
	for _, option := range options {
		if l == option {
			return true
		}
	}
	return false
}
