package clt

import (
	"fmt"
	"strings"
)

// ValidationFunc is a type alias for a validator that takes a string and
// returns true if it passes validation.  Error provides a helpful error
// message to the user that is shown before re-asking the question.
type ValidationFunc func(s string) (bool, error)

// ValidateYesNo is a validation function that ensures that a yes/no question
// receives a valid response.  Use IsYes or IsNo to test for a particular
// yes or no response.
func ValidateYesNo() ValidationFunc {
	return func(s string) (bool, error) {
		options := []string{"yes", "y", "no", "n"}
		return validateOptions(strings.ToLower(s), options)
	}
}

// IsYes returns true if the response is some form of yes or y
func IsYes(s string) bool {
	options := []string{"yes", "y"}
	ok, _ := validateOptions(strings.ToLower(s), options)
	return ok
}

// IsNo returns true if the response is some form of no or n
func IsNo(s string) bool {
	options := []string{"no", "n"}
	ok, _ := validateOptions(strings.ToLower(s), options)
	return ok
}

// AllowedOptions builds a new validator from a []string of options
func AllowedOptions(options []string) ValidationFunc {
	return func(s string) (bool, error) {
		return validateOptions(strings.ToLower(s), options)
	}
}

// Required validates that the length of the input is greater than 0
func Required() ValidationFunc {
	return func(s string) (bool, error) {
		switch {
		case len(s) > 0:
			return true, nil
		default:
			return false, fmt.Errorf("A response is required.")
		}
	}
}

// validateOptions returns true if the given string appears in the list of valid
// options.
func validateOptions(s string, options []string) (bool, error) {
	l := strings.ToLower(s)
	for _, option := range options {
		if l == option {
			return true, nil
		}
	}
	return false, fmt.Errorf("%s is a not a valid option. Valid options are %v", s, options)
}
