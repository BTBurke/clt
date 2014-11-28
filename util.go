package clt

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// StringBuilder represents a base class for building strings
type StringBuilder struct {
	msg bytes.Buffer
}

// Add concatenates a new string to existing message
func (s *StringBuilder) Add(format string, args ...interface{}) {
	s.msg.WriteString(fmt.Sprintf(format, args...))
}

// Addln concatenates a new string with an newline ending if one does not
// already exist
func (s *StringBuilder) Addln(format string, args ...interface{}) {
	if strings.HasSuffix(format, "\n") {
		s.Add(format, args...)
	} else {
		s.Add(format+"\n", args...)
	}
}

// NewLine concatenates a newline character to the existing string
func (s *StringBuilder) NewLine() {
	s.msg.WriteString("\n")
}

// NewLines concatenates multiple newline characters to the existing string
func (s *StringBuilder) NewLines(num int) {
	for i := 0; i < num; i++ {
		s.msg.WriteString("\n")
	}
}

// Render prints the string to StdOut
func (s *StringBuilder) Render() {
	fmt.Printf(s.Finalize())
}

// Finalize returns the message as a string
func (s *StringBuilder) Finalize() string {
	return s.msg.String()
}

// NewStringBuilder returns a pointed to an empty StringBuilder struct
func NewStringBuilder() *StringBuilder {
	return &StringBuilder{}
}

// ConsoleInput represents an interactive console session used to prompt the user
// for input.  Prompt is rendered right before the input, with an optional Default
// that will be returned if nothing is entered.
type ConsoleInput struct {
	Prompt   string
	Response string
	Default  string
}

// Ask renders the prompt, a default if it exists, and collects the response
func (c *ConsoleInput) Ask() {
	if len(c.Default) > 0 {
		fmt.Printf("%s  [%s]: ", c.Prompt, c.Default)
	} else {
		fmt.Printf("%s: ", c.Prompt)
	}
	reader := bufio.NewReader(os.Stdin)
	c.Response, _ = reader.ReadString('\n')
}

// Ask renders the prompt and default if it exists, collects the response and
// then validates it using ValidationFunc.  If the response does not pass validation,
// an error message is shown and ask is called again.
func (c *ConsoleInput) AskValidate(f ValidationFunc) {
	c.Ask()
	if !c.Validate(f) {
		fmt.Printf("Error: %s is not a valid response.", c.Response)
		c.AskValidate(f)
	}
}

// Get renders the prompt using Ask() but explicitly returns the response
// as a string.
func (c *ConsoleInput) Get() string {
	c.Ask()
	return c.Response
}

// Validate uses ValidationFunc which validates the response, returning a bool.
func (c *ConsoleInput) Validate(f ValidationFunc) bool {
	return f(c.Response)
}

// NewConsoleInput returns a pointer to a new ConsoleInput
func NewConsoleInput() *ConsoleInput {
	return &ConsoleInput{}
}
