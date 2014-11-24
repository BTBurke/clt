package mcclintock

import (
	"fmt"
)

type color struct {
	before int
	after int
}

func (c color) Codes() (int, int) { return c.before, c.after }

type textstyle struct {
	before int
	after int
}

func (t textstyle) Codes() (int, int) { return t.before, t.after }

type Style interface {
	Codes() (int, int)
}

var (
	// Colors
	Black   = color{0, 9}
	Red     = color{1, 9}
	Green   = color{2, 9}
	Yellow  = color{3, 9}
	Blue    = color{4, 9}
	Magenta = color{5, 9}
	Cyan    = color{6, 9}
	White   = color{7, 9}
	Default = color{9, 9}

	// Shortcut Colors
	K = color{0, 9}
	R = color{1, 9}
	G = color{2, 9}
	Y = color{3, 9}
	B = color{4, 9}
	M = color{5, 9}
	C = color{6, 9}
	W = color{7, 9}
	D = color{9, 9}

	// Textstyles
	Bold          = textstyle{1, 22}
	Italic        = textstyle{3, 23}
	Underline     = textstyle{4, 24}	

)

type text struct {}

func Text() *text {
	return &text{}
}

func (t *text) Red(s string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[39m")
}