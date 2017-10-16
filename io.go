package clt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// OutputWriter represents a base class for building strings
type OutputWriter struct {
	msg    bytes.Buffer
	writer io.Writer
}

// Add concatenates a new string to existing message
func (s *OutputWriter) Add(format string, args ...interface{}) *OutputWriter {
	s.msg.WriteString(fmt.Sprintf(format, args...))
	return s
}

// Addln concatenates a new string with an newline ending if one does not
// already exist
func (s *OutputWriter) Addln(format string, args ...interface{}) *OutputWriter {
	if strings.HasSuffix(format, "\n") {
		s.Add(format, args...)
	} else {
		s.Add(format+"\n", args...)
	}
	return s
}

// NewLine concatenates a newline character to the existing string
func (s *OutputWriter) NewLine() *OutputWriter {
	s.msg.WriteString("\n")
	return s
}

// NewLines concatenates multiple newline characters to the existing string
func (s *OutputWriter) NewLines(num int) *OutputWriter {
	for i := 0; i < num; i++ {
		s.msg.WriteString("\n")
	}
	return s
}

// Render prints the string to StdOut
func (s *OutputWriter) Render() {
	if s.writer == nil {
		s.writer = os.Stdout
	}
	fmt.Fprintf(s.writer, s.Finalize())
}

// Finalize returns the message as a string
func (s *OutputWriter) Finalize() string {
	return s.msg.String()
}

// NewOutputWriter returns a pointed to an empty OutputWriter struct
func NewOutputWriter() *OutputWriter {
	return &OutputWriter{
		writer: os.Stdout,
	}
}

// InputReader represents an interactive console session used to prompt the user
// for input.  Prompt is rendered right before the input, with an optional Default
// that will be returned if nothing is entered.
type InputReader struct {
	Prompt   string
	Response string
	Default  string
	ValHint  string

	writer *bufio.Writer
	reader *bufio.Reader
}

// Ask renders the prompt, a default and validation hint if it exists, and collects the response
func (c *InputReader) Ask() (err error) {
	if c.writer == nil {
		c.writer = bufio.NewWriter(os.Stdout)
	}
	switch {
	case len(c.Default) > 0:
		fmt.Fprintf(c.writer, "%s  [%s]: ", c.Prompt, c.Default)
	case len(c.ValHint) > 0:
		fmt.Fprintf(c.writer, "%s (%s): ", c.Prompt, c.ValHint)
	default:
		fmt.Fprintf(c.writer, "%s: ", c.Prompt)
	}
	c.writer.Flush()

	c.Response, err = c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	if len(c.Default) > 0 && len(c.Response) == 0 {
		c.Response = c.Default
	}
	c.Response = strings.TrimSpace(c.Response)
	return nil
}

// NewInputReader returns a pointer to a new InputReader
func NewInputReader() *InputReader {
	return &InputReader{
		reader: bufio.NewReader(os.Stdin),
		writer: bufio.NewWriter(os.Stdout),
	}
}
