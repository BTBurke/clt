package clt

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	success int = iota
	fail
)

const (
	spinner int = iota
	bar
)

type Progress struct {
	// Prompt to display before spinner or bar
	Prompt string
	// Approximate length of the total progress display, including
	// the prompt and the ..., does not include status indicator
	// at the end (e.g, the spinner, FAIL, OK, or XX%)
	DisplayLength int
	// Other non-exported fields
	style int
	cf    chan float64
	c     chan int
}

// NewProgressSpinner returns a new spinner with prompt <message>
// display length defaults to 30.
func NewProgressSpinner(format string, args ...interface{}) *Progress {
	return &Progress{
		style:         spinner,
		Prompt:        fmt.Sprintf(format, args...),
		DisplayLength: 30,
	}
}

// NewProgressBar returns a new progress bar with prompt <message>
// display length defaults to 20
func NewProgressBar(format string, args ...interface{}) *Progress {
	return &Progress{
		style:         bar,
		Prompt:        fmt.Sprintf(format, args...),
		DisplayLength: 20,
	}
}

// Start launches a Goroutine to render the progress bar or spinner
// and returns control to the caller for further processing.  Spinner
// will update automatically every 250ms until Success() or Fail() is
// called.  Bars will update by calling Update(<pct_complete>).  You
// must always finally call either Success() or Fail() to terminate
// the go routine.
func (p *Progress) Start() {
	switch p.style {
	case spinner:
		p.c = make(chan int)
		go renderSpinner(p, p.c)
	case bar:
		p.cf = make(chan float64)
		go renderBar(p, p.cf)
	}
}

// Success should be called on a progress bar or spinner
// after completion is successful
func (p *Progress) Success() {
	switch p.style {
	case spinner:
		p.c <- success
	case bar:
		p.cf <- -1.0
	}
}

// Fail should be called on a progress bar or spinner
// if a failure occurs
func (p *Progress) Fail() {
	switch p.style {
	case spinner:
		p.c <- fail
	case bar:
		p.cf <- -2.0
	}
}

func renderSpinner(p *Progress, c chan int) {
	promptLen := len(p.Prompt)
	dotLen := p.DisplayLength - promptLen
	if dotLen < 3 {
		dotLen = 3
	}
	for i := 0; ; i++ {
		select {
		case result := <-c:
			switch result {
			case success:
				fmt.Printf("\x1b[?25h\r%s%s[%s]\n", p.Prompt, strings.Repeat(".", dotLen), Style(Green).ApplyTo("OK"))
			case fail:
				fmt.Printf("\x1b[?25h\r%s%s[%s]\n", p.Prompt, strings.Repeat(".", dotLen), Style(Red).ApplyTo("FAIL"))
			}
			return
		default:
			fmt.Printf("\x1b[?25l\r%s%s[%s]", p.Prompt, strings.Repeat(".", dotLen), spinLookup(i))
			time.Sleep(time.Duration(250) * time.Millisecond)
		}
	}
}

func spinLookup(i int) string {
	switch int(math.Mod(float64(i), 4.0)) {
	case 0:
		return "|"
	case 1:
		return "/"
	case 2:
		return "-"
	case 3:
		return "\\"
	}
	return ""
}

func renderBar(p *Progress, c chan float64) {
	var result float64
	eqLen := 0
	spLen := p.DisplayLength

	for {
		select {
		case result = <-c:
			eqLen = int(result * float64(p.DisplayLength))
			spLen = p.DisplayLength - eqLen
			switch {
			case result == -1.0:
				fmt.Printf("\x1b[?25l\r%s: [%s] %s", p.Prompt, strings.Repeat("=", p.DisplayLength), Style(Green).ApplyTo("100%"))
				fmt.Printf("\x1b[?25h\n")
				return
			case result == -2.0:
				fmt.Printf("\x1b[?25l\r%s: [%s] %s", p.Prompt, strings.Repeat("X", p.DisplayLength), Style(Red).ApplyTo("FAIL"))
				fmt.Printf("\x1b[?25h\n")
				return
			case result >= 0.0:
				fmt.Printf("\x1b[?25l\r%s: [%s%s] %2.0f%%", p.Prompt, strings.Repeat("=", eqLen), strings.Repeat(" ", spLen), 100.0*result)
			}
		default:
			fmt.Printf("\x1b[?25l\r%s: [%s%s] %2.0f%%", p.Prompt, strings.Repeat("=", eqLen), strings.Repeat(" ", spLen), 100.0*result)
		}
	}
}

// Progress updates the progress bar using a number [0, 1.0] to represent
// the percentage complete
func (p *Progress) Update(pct float64) {
	if pct >= 1.0 {
		pct = 1.0
	}
	p.cf <- pct
}
