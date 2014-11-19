package mcclintock

import (
	"fmt"
	"testing"
)

func TestTerminalSizeCheck(t *testing.T) {
	h, w, err := getTerminalSize()
	if err != nil || h == -1 || w == -1 {
		fmt.Printf("Cannot determine terminal size for Table.  McClintock will still work, but will not be able to automagically determine sizes.")
	}
}
