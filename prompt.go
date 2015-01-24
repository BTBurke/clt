package clt

import "os"

// Say is a thin wrapper around fmt.Printf.  You can use either, but Say
// can be used to highlight code meant specifically for user interaction.
func Say(format string, args ...interface{}) {
	s := buildSay(format, args...)
	s.Render()
}

func buildSay(format string, args ...interface{}) *StringBuilder {
	s := NewStringBuilder()
	s.Addln(format, args...)
	return s
}

// Pause will stop program execution until the user presses enter
func Pause() {
	c := NewConsoleInput()
	c.Prompt = "\n\nPress [Enter] to continue.\n\n"
	c.Ask()
}

// SayAndPause will print the message to Stdout then wait for the user
// to press [enter] to continue
func SayAndPause(format string, args ...interface{}) {
	s := buildSay(format, args...)
	s.Render()
	Pause()
}

// Warn gives an informational warning message to the user in format
// Warning: <user defined string>
func Warn(format string, args ...interface{}) {
	s := NewStringBuilder()
	s.Add("%s: ", Style(Yellow).ApplyTo("Warning"))
	s.Add(format, args...)
	s.Render()
}

// Error gives an informational error message to the user in format
// Error: <user defined string>.  Exits the program returning status
// code 1
func Error(format string, args ...interface{}) {
	s := NewStringBuilder()
	s.Add("%s: ", Style(Red).ApplyTo("Error:"))
	s.Add(format, args...)
	s.Render()
	os.Exit(1)
}

// Info gives an informational message to the user in format
// Info: <user defined string>
func Info(format string, args ...interface{}) {
	s := NewStringBuilder()
	s.Add("Info: ")
	s.Add(format, args...)
	s.Render()
}
