package clt

import (
	"testing"
)

// TestYesNo tests inputs using yes and no values
func TestYesNo(t *testing.T) {
	cases := []string{"yes", "y", "NO", "n"}

	for _, case1 := range cases {
		if !ValidateYesNo(case1) {
			t.Errorf("Failed YesNoValidation for %s", case1)
		}
	}

	cases = []string{"yes", "YES", "Y", "y"}
	for _, y := range cases {
		if !IsYes(y) {
			t.Errorf("Failed IsYes validation for %s", y)
		}
	}

	cases = []string{"no", "NO", "N", "n"}
	for _, n := range cases {
		if !IsNo(n) {
			t.Errorf("Failed IsNo validation for %s", n)
		}
	}
}

// TestOptionValidator tests building a new validator function with a list
// of options.
func TestOptionValidator(t *testing.T) {
	passCases := []string{"go", "capitals"}
	failCases := []string{"bruins", "rangers"}
	options := []string{"go", "washington", "capitals"}

	valFunc := OptionValidator(options)
	for _, p := range passCases {
		if !ValidateWithFunc(p, valFunc) {
			t.Errorf("Failed OptionValidation for %s", p)
		}
	}

	for _, f := range failCases {
		if ValidateWithFunc(f, valFunc) {
			t.Errorf("Failed OptionValidation for %s", f)
		}
	}

}
