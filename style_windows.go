// +build windows

package clt

import (
	"bytes"
)

// Color represents a ANSI-coded color style for text
type Color struct {
	before int
	after  int
}

// Codes returns ANSI styling values for a color
func (c Color) Codes() (int, int) { return c.before, c.after }

// Textstyle represents a ANSI-coded text style
type Textstyle struct {
	before int
	after  int
}

// Codes returns ANSI styling values for a textstyle
func (t Textstyle) Codes() (int, int) { return t.before, t.after }

// Styler is an interface that is fulfilled by either a Color
// or Textstyle to be applied to a string
type Styler interface {
	Codes() (int, int)
}

// Style represents a computed style from one or more colors or textstyles
// as the ANSI code suitable for terminal output
type Style struct {
	before string
	after  string
}

// ApplyTo applies styles created using the Styled command to a string
// to generate an styled output using ANSI terminal codes
func (s *Style) ApplyTo(content string) string {
	var out bytes.Buffer
	out.WriteString(s.before)
	out.WriteString(content)
	out.WriteString(s.after)
	return out.String()
}

var (
	// Colors
	Black   = Color{30, 39}
	Red     = Color{31, 39}
	Green   = Color{32, 39}
	Yellow  = Color{33, 39}
	Blue    = Color{34, 39}
	Magenta = Color{35, 39}
	Cyan    = Color{36, 39}
	White   = Color{37, 39}
	Default = Color{39, 39}

	// Shortcut Colors
	K   = Color{30, 39}
	R   = Color{31, 39}
	G   = Color{32, 39}
	Y   = Color{33, 39}
	B   = Color{34, 39}
	M   = Color{35, 39}
	C   = Color{36, 39}
	W   = Color{37, 39}
	Def = Color{39, 39}

	// Textstyles
	Bold      = Textstyle{1, 22}
	Italic    = Textstyle{3, 23}
	Underline = Textstyle{4, 24}
)

// Background returns a style that sets the background to the appropriate color
func Background(c Color) Color {
	c.before += 10
	c.after += 10
	return c
}

// Styled contructs a composite style from one of more color or textstyle values.  Styles
// can be applied to a string via ApplyTo or as a shortcut use SStyled which returns a string directly
// Example:  Styled(White, Underline)
func Styled(s ...Styler) *Style {
	switch {
	case len(s) == 1:
		var computedStyle Style
		computedStyle.before = ""
		computedStyle.after = ""
		return &computedStyle
	case len(s) > 1:
		var computedStyle Style
		var beforeConcat, afterConcat bytes.Buffer

		beforeConcat.WriteString("")
		afterConcat.WriteString("")

		for idx := range s {
			if idx < len(s)-1 {
				beforeConcat.WriteString("")
				afterConcat.WriteString("")
			} else {
				beforeConcat.WriteString("")
				afterConcat.WriteString("")
			}
		}
		computedStyle.before = beforeConcat.String()
		computedStyle.after = afterConcat.String()
		return &computedStyle
	}
	return &Style{}
}

// SStyled is a shorter version of Styled(s...).ApplyTo(content)
func SStyled(content string, s ...Styler) string {
	return Styled(s...).ApplyTo(content)
}
