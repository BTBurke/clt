package mcclintock

import (
	"bytes"
	"fmt"
)

type color struct {
	before int
	after  int
}

func (c color) codes() (int, int) { return c.before, c.after }

type textstyle struct {
	before int
	after  int
}

func (t textstyle) codes() (int, int) { return t.before, t.after }

type styleInterface interface {
	codes() (int, int)
}

// style represents a computed style from one or more colors or textstyles
// as the ANSI code suitable for terminal output
type style struct {
	before string
	after  string
}

// ApplyTo applies styles created using the Style command to a string
// to generate an styled output using ANSI terminal codes
func (s *style) ApplyTo(content string) string {
	var out bytes.Buffer
	out.WriteString(s.before)
	out.WriteString(content)
	out.WriteString(s.after)
	return out.String()
}

var (
	// Colors
	Black   = color{30, 39}
	Red     = color{31, 39}
	Green   = color{32, 39}
	Yellow  = color{33, 39}
	Blue    = color{34, 39}
	Magenta = color{35, 39}
	Cyan    = color{36, 39}
	White   = color{37, 39}
	Default = color{39, 39}

	// Shortcut Colors
	K   = color{30, 39}
	R   = color{31, 39}
	G   = color{32, 39}
	Y   = color{33, 39}
	B   = color{34, 39}
	M   = color{35, 39}
	C   = color{36, 39}
	W   = color{37, 39}
	Def = color{39, 39}

	// Textstyles
	Bold      = textstyle{1, 22}
	Italic    = textstyle{3, 23}
	Underline = textstyle{4, 24}
)

// Background returns a style that sets the background to the appropriate color
func Background(c *color) *color {
	c.before += 10
	c.after += 10
	return c
}

func Style(s ...styleInterface) *style {
	switch {
	case len(s) == 1:
		bef, aft := s[0].codes()
		var computedStyle style
		computedStyle.before = fmt.Sprintf("\x1b[%vm", bef)
		computedStyle.after = fmt.Sprintf("\x1b[%vm", aft)
		return &computedStyle
	case len(s) > 1:
		var computedStyle style
		var beforeConcat, afterConcat bytes.Buffer

		beforeConcat.WriteString("\x1b[")
		afterConcat.WriteString("\x1b[")

		var bef, aft int
		for idx, sty := range s {
			bef, aft = sty.codes()
			if idx < len(s)-1 {
				beforeConcat.WriteString(fmt.Sprintf("%v; ", bef))
				afterConcat.WriteString(fmt.Sprintf("%v; ", aft))
			} else {
				beforeConcat.WriteString(fmt.Sprintf("%vm", bef))
				afterConcat.WriteString(fmt.Sprintf("%vm", aft))
			}
		}
		computedStyle.before = beforeConcat.String()
		computedStyle.after = afterConcat.String()
		return &computedStyle
	}
	return &style{}
}
