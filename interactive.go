package clt

import (
	"bufio"
	"fmt"
	"io"

	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// InteractiveSession creates a system for collecting user input
// in response to questions and choices
type InteractiveSession struct {
	Prompt   string
	Response string
	Default  string
	ValHint  string

	input  *bufio.Reader
	output io.Writer
}

// NewInteractiveSession returns a new InteractiveSession outputting to Stdout
// and reading from Stdin
func NewInteractiveSession() *InteractiveSession {
	return &InteractiveSession{
		input:  bufio.NewReader(os.Stdin),
		output: bufio.NewWriter(os.Stdout),
	}
}

// Say
func (i *InteractiveSession) Say(format string, args ...interface{}) *InteractiveSession {
	fmt.Fprintf(i.output, format, args...)
	return i
}

// Pause is a terminator that will render long-form text added via the another method
// that returns *InteractiveSession and will wait for the user to press enter to continue.
// It is useful for long-form content or paging.
func (i *InteractiveSession) Pause() {
	i.Prompt = "\n\nPress [Enter] to continue.\n\n"
	i.get()
}

func (i *InteractiveSession) get() (err error) {
	switch {
	case len(i.Default) > 0:
		fmt.Fprintf(i.output, "%s  [%s]: ", i.Prompt, i.Default)
	case len(i.ValHint) > 0:
		fmt.Fprintf(i.output, "%s (%s): ", i.Prompt, i.ValHint)
	default:
		fmt.Fprintf(i.output, "%s: ", i.Prompt)
	}

	i.Response, err = i.input.ReadString('\n')
	if err != nil {
		return err
	}
	if len(i.Default) > 0 && len(i.Response) == 0 {
		i.Response = i.Default
	}
	i.Response = strings.TrimSpace(i.Response)
	return nil
}

// Warn adds an informational warning message to the user in format
// Warning: <user defined string>
func (i *InteractiveSession) Warn(format string, args ...interface{}) *InteractiveSession {
	fmt.Fprintf(i.output, "%s: %s", Style(Yellow).ApplyTo("Warning"), fmt.Sprintf(format, args...))
	return i
}

// Error is a terminator that gives an informational error message to the user in format
// Error: <user defined string>.  Exits the program returning status code 1
func (i *InteractiveSession) Error(format string, args ...interface{}) {
	fmt.Fprintf(i.output, "\n\n%s: %s\n", Style(Red).ApplyTo("Error:"), fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Ask is a terminator for an interactive session that results in returning the user
// input.  Validators can optionally be applied to ensure that acceptable input is returned
// or the question will be asked again.
func (i *InteractiveSession) Ask(prompt string, validators ...ValidationFunc) string {
	return i.ask(prompt, "", "", validators...)
}

func (i *InteractiveSession) AskWithDefault(prompt string, defaultChoice string, validators ...ValidationFunc) string {
	return i.ask(prompt, defaultChoice, "", validators...)
}

func (i *InteractiveSession) AskWithHint(prompt string, hint string, validators ...ValidationFunc) string {
	return i.ask(prompt, "", hint, validators...)
}

func (i *InteractiveSession) ask(prompt string, def string, hint string, validators ...ValidationFunc) string {
	i.Prompt = prompt
	i.Default = def
	i.ValHint = hint
	i.get()
	for _, validator := range validators {
		if ok, err := validator(i.Response); !ok {
			i.Say("\nError: %s\n\n", err)
			i.ask(prompt, def, hint, validators...)
		}
	}
	return i.Response
}

// AskPassword is a terminator that asks for a password and does not echo input
// to the terminal.
func (i *InteractiveSession) AskPassword(validators ...ValidationFunc) string {
	rw := bufio.NewReadWriter(i.input, bufio.NewWriter(i.output))
	term := terminal.NewTerminal(rw, "")
	pw, err := term.ReadPassword("Password: ")
	if err != nil {
		i.Error("\n%s\n\n", err)
	}
	for _, validator := range validators {
		if ok, err := validator(strings.TrimSpace(pw)); !ok {
			i.Say("\nError: %s\n\n", err)
			i.AskPassword(validators...)
		}
	}
	return strings.TrimSpace(pw)
}

func (i *InteractiveSession) YesNo(prompt string, defaultChoice string) string {
	switch def := strings.ToLower(defaultChoice); def {
	case "y", "yes":
		i.Default = "Y/n"
	default:
		i.Default = "y/N"
	}
	return i.ask(prompt, i.Default, "", ValidateYesNo())
}

func (i *InteractiveSession) AskFromTable(prompt string, choices map[string]string, def string) string {
	t := NewTable(2)
	var allKeys []string
	for key, choice := range choices {
		t.AddRow(key, choice)
		allKeys = append(allKeys, key)
	}
	tAsString := t.renderTableAsString()

	i.Prompt = fmt.Sprintf("%s\n%s\n\nChoice [%s]: ", prompt, tAsString, choices[def])
	i.Default = def
	i.get()
	if ok, err := AllowedOptions(allKeys)(i.Response); !ok {
		i.Say("\nError: %s\n\n", err)
		i.AskFromTable(prompt, choices, def)
	}

	return strings.TrimSpace(i.Response)
}
