package clt

import (
	"bufio"
	"fmt"
	"io"
	"sort"

	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

type getOpt int

const (
	noColon getOpt = iota
)

// InteractiveSession creates a system for collecting user input
// in response to questions and choices
type InteractiveSession struct {
	Prompt  string
	Default string
	ValHint string

	response string
	input    *bufio.Reader
	output   io.Writer
}

// NewInteractiveSession returns a new InteractiveSession outputting to Stdout
// and reading from Stdin by default, but other inputs and outputs may be specified
// with SessionOptions
func NewInteractiveSession(opts ...SessionOption) *InteractiveSession {
	i := &InteractiveSession{
		input:  bufio.NewReader(os.Stdin),
		output: os.Stdout,
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

// SessionOption optionally configures aspects of the interactive session
type SessionOption func(i *InteractiveSession)

// WithInput uses an input other than os.Stdin
func WithInput(r io.Reader) SessionOption {
	return func(i *InteractiveSession) {
		i.input = bufio.NewReader(r)
	}
}

// WithOutput uses an output other than os.Stdout
func WithOutput(w io.Writer) SessionOption {
	return func(i *InteractiveSession) {
		i.output = w
	}
}

// Reset allows reuse of the same interactive session by reseting its state and keeping
// its current input and output
func (i *InteractiveSession) Reset() {
	i.Prompt = ""
	i.Default = ""
	i.ValHint = ""
	i.response = ""
}

// Say is a short form of fmt.Fprintf but allows you to chain additional terminators to
// the interactive session to collect user input
func (i *InteractiveSession) Say(format string, args ...interface{}) *InteractiveSession {
	fmt.Fprintf(i.output, fmt.Sprintf("\n%s\n", format), args...)
	return i
}

// Pause is a terminator that will render long-form text added via the another method
// that returns *InteractiveSession and will wait for the user to press enter to continue.
// It is useful for long-form content or paging.
func (i *InteractiveSession) Pause() {
	i.Prompt = "\nPress [Enter] to continue."
	i.get(noColon)
}

// PauseWithPrompt is a terminator that will render long-form text added via the another method
// that returns *InteractiveSession and will wait for the user to press enter to continue.
// This will use the custom prompt specified by format and args.
func (i *InteractiveSession) PauseWithPrompt(format string, args ...interface{}) {
	i.Prompt = fmt.Sprintf(format, args...)
	i.get(noColon)
}

func (i *InteractiveSession) get(opts ...getOpt) (err error) {
	contains := func(wanted getOpt) bool {
		for _, opt := range opts {
			if opt == wanted {
				return true
			}
		}
		return false
	}

	if i.output == nil {
		i.output = bufio.NewWriter(os.Stdout)
	}
	if i.input == nil {
		i.input = bufio.NewReader(os.Stdin)
	}

	switch {
	case len(i.Default) > 0:
		fmt.Fprintf(i.output, "%s  [%s]: ", i.Prompt, i.Default)
	case len(i.ValHint) > 0:
		fmt.Fprintf(i.output, "%s (%s): ", i.Prompt, i.ValHint)
	case contains(noColon):
		fmt.Fprintf(i.output, "%s", i.Prompt)
	case len(i.Prompt) > 0:
		fmt.Fprintf(i.output, "%s: ", i.Prompt)
	default:
	}

	i.response, err = i.input.ReadString('\n')
	if err != nil {
		return err
	}
	i.response = strings.TrimRight(i.response, " \n\r")
	if len(i.Default) > 0 && len(i.response) == 0 {
		switch i.Default {
		case "y/N":
			i.response = "n"
		case "Y/n":
			i.response = "y"
		default:
			i.response = i.Default
		}
	}

	return nil
}

// Warn adds an informational warning message to the user in format
// Warning: <user defined string>
func (i *InteractiveSession) Warn(format string, args ...interface{}) *InteractiveSession {
	fmt.Fprintf(i.output, "\n%s: %s\n", Styled(Yellow).ApplyTo("Warning"), fmt.Sprintf(format, args...))
	return i
}

// Error is a terminator that gives an informational error message to the user in format
// Error: <user defined string>.  Exits the program returning status code 1
func (i *InteractiveSession) Error(format string, args ...interface{}) {
	fmt.Fprintf(i.output, "\n\n%s: %s\n", Styled(Red).ApplyTo("Error:"), fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Ask is a terminator for an interactive session that results in returning the user's
// input.  Validators can optionally be applied to ensure that acceptable input is returned
// or the question will be asked again.
func (i *InteractiveSession) Ask(prompt string, validators ...ValidationFunc) string {
	return i.ask(prompt, "", "", validators...)
}

// AskWithDefault is like ask, but sets a default choice that the user can select by pressing enter.
func (i *InteractiveSession) AskWithDefault(prompt string, defaultChoice string, validators ...ValidationFunc) string {
	return i.ask(prompt, defaultChoice, "", validators...)
}

// AskWithHint is like ask, but gives a hint about the proper format of the response.  This is useful
// combined with a validation function to get input in the right format.
func (i *InteractiveSession) AskWithHint(prompt string, hint string, validators ...ValidationFunc) string {
	return i.ask(prompt, "", hint, validators...)
}

func (i *InteractiveSession) ask(prompt string, def string, hint string, validators ...ValidationFunc) string {
	i.Prompt = prompt
	i.Default = def
	i.ValHint = hint
	i.get()
	for _, validator := range validators {
		if ok, err := validator(i.response); !ok {
			i.Say("\nError: %s\n\n", err)
			i.ask(prompt, def, hint, validators...)
		}
	}
	return i.response
}

// AskPassword is a terminator that asks for a password and does not echo input
// to the terminal.
func (i *InteractiveSession) AskPassword(validators ...ValidationFunc) string {
	return askPassword(i, "Password: ", validators...)
}

// AskPasswordPrompt is a terminator that asks for a password with a custom prompt
func (i *InteractiveSession) AskPasswordPrompt(prompt string, validators ...ValidationFunc) string {
	return askPassword(i, prompt, validators...)
}

func askPassword(i *InteractiveSession, prompt string, validators ...ValidationFunc) string {
	fmt.Fprintf(i.output, "Password: ")
	pw, err := terminal.ReadPassword(0)
	if err != nil {
		i.Error("\n%s\n", err)
	}

	pwS := strings.TrimSpace(string(pw))
	for _, validator := range validators {
		if ok, err := validator(pwS); !ok {
			i.Say("\nError: %s\n\n", err)
			i.AskPassword(validators...)
		}
	}
	return pwS
}

// AskYesNo asks the user a yes or no question with a default value.  Defaults of `y` or `yes` will
// set the default to yes.  Anything else will default to no.  You can use IsYes or IsNo to act on the response
// without worrying about what version of y, Y, YES, yes, etc. that the user entered.
func (i *InteractiveSession) AskYesNo(prompt string, defaultChoice string) string {
	switch def := strings.ToLower(defaultChoice); def {
	case "y", "yes":
		i.Default = "Y/n"
	default:
		i.Default = "y/N"
	}
	return i.ask(prompt, i.Default, "", ValidateYesNo())
}

// AskFromTable creates a table to select choices from.  It has a built-in validation function that will
// ensure that only the options listed are valid choices.
func (i *InteractiveSession) AskFromTable(prompt string, choices map[string]string, def string) string {
	t := NewTable(2).
		ColumnHeaders("Option", "")
	var allKeys []string
	for key := range choices {
		allKeys = append(allKeys, key)
	}
	sort.Strings(allKeys)

	for _, key := range allKeys {
		t.AddRow(key, choices[key])
	}
	tAsString := t.AsString()

	i.Prompt = fmt.Sprintf("\n%s%s\nChoice", prompt, tAsString)
	i.Default = def
	i.get()
	if ok, err := AllowedOptions(allKeys)(i.response); !ok {
		i.Say("\nError: %s\n\n", err)
		i.AskFromTable(prompt, choices, def)
	}

	return strings.TrimSpace(i.response)
}
