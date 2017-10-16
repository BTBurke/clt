package clt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// InteractiveSession creates a system for collecting user input
// in response to questions and choices
type InteractiveSession struct {
	input  *InputReader
	output *OutputWriter
}

// NewInteractiveSession returns a new InteractiveSession outputting to Stdout
// and reading from Stdin
func NewInteractiveSession() *InteractiveSession {
	return &InteractiveSession{
		input:  NewInputReader(),
		output: NewOutputWriter(),
	}
}

// SayAndThen creates a long-form prompt and then asks for some type of user input.  You can use this
// to add information ahead of asking for user input.  It can be used with any terminator such as Render,
// Pause, Ask, or Choice.
func (i *InteractiveSession) SayAndThen(format string, args ...interface{}) *InteractiveSession {
	i.output.Add(format, args...)
	return i
}

// Say is a terminator that immediately renders information to the console and does not ask for
// any type of user input.
func (i *InteractiveSession) Say(format string, args ...interface{}) {
	i.SayAndThen(format, args...).Render()
}

// Pause is a terminator that will render long-form text added via the another method
// that returns *InteractiveSession and will wait for the user to press enter to continue.
// It is useful for long-form content or paging.
func (i *InteractiveSession) Pause() {
	c := NewInputReader()
	c.Prompt = i.output.Add("\n\nPress [Enter] to continue.\n\n").Finalize()
	i.input.Ask()
}

// Warn adds an informational warning message to the user in format
// Warning: <user defined string>
func (i *InteractiveSession) Warn(format string, args ...interface{}) *InteractiveSession {
	i.output.Add("%s: %s", Style(Yellow).ApplyTo("Warning"), fmt.Sprintf(format, args...))
	return i
}

// Error is a terminator that gives an informational error message to the user in format
// Error: <user defined string>.  Exits the program returning status code 1
func (i *InteractiveSession) Error(format string, args ...interface{}) {
	i.output.Add("%s: %s", Style(Red).ApplyTo("Error:"), fmt.Sprintf(format, args...)).Render()
	os.Exit(1)
}

// Render is a terminator for an interactive session that provides information to the
// user but does not ask for any input.  To retrieve info, you should terminate
// with Ask, Choice, or another terminator that returns user input.
func (i *InteractiveSession) Render() {
	i.output.Render()
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
	c := NewInputReader()
	c.Prompt = prompt
	i.input.Default = def
	i.input.ValHint = hint
	i.Render()
	i.input.Ask()
	for _, validator := range validators {
		if ok, err := validator(i.input.Response); !ok {
			i.Say("Error: %s", err)
			i.Ask(prompt, validators...)
		}
	}
	return i.input.Response
}

// AskPassword is a terminator that asks for a password and does not echo input
// to the terminal.
func (i *InteractiveSession) AskPassword(validators ...ValidationFunc) string {
	rw := bufio.NewReadWriter(i.input.reader, bufio.NewWriter(i.output.writer))
	term := terminal.NewTerminal(rw, "")
	pw, err := term.ReadPassword("Password: ")
	if err != nil {
		i.Error("%s", err)
	}
	for _, validator := range validators {
		if ok, err := validator(strings.TrimSpace(pw)); !ok {
			i.Say("Error: %s", err)
			i.AskPassword(validators...)
		}
	}
	return strings.TrimSpace(pw)
}

func (i *InteractiveSession) AskYesNo(prompt string, defaultChoice string) string {
	switch def := strings.ToLower(defaultChoice); def {
	case "y", "yes":
		i.input.Default = "Y/n"
	default:
		i.input.Default = "y/N"
	}
	return i.ask(prompt, i.input.Default, "", ValidateYesNo())
}

func (i *InteractiveSession) AskFromTable(prompt string, choices map[string]string, def string) string {
	t := NewTable(2)
	var allKeys []string
	for key, choice := range choices {
		t.AddRow(key, choice)
		allKeys = append(allKeys, key)
	}
	tAsString := t.renderTableAsString()

	i.input.Prompt = fmt.Sprintf("%s\n%s\nChoice [%s]: ", prompt, tAsString, choices[def])
	i.input.Default = def
	i.Render()
	i.input.Ask()
	if ok, err := AllowedOptions(allKeys)(i.input.Response); !ok {
		i.Say("Error: %s", err)
		i.AskFromTable(prompt, choices, def)
	}

	return strings.TrimSpace(i.input.Response)
}
