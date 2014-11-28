package clt

func Say(format string, args ...interface{}) {
	s := buildSay(format, args...)
	s.Render()
}

func buildSay(format string, args ...interface{}) *StringBuilder {
	s := NewStringBuilder()
	s.Addln(format, args...)
	return s
}

func Pause() {
	c := NewConsoleInput()
	c.Prompt = "\n\nPress [Enter} to continue.\n\n"
	c.Ask()
}

func SayAndPause(format string, args ...interface{}) {
	s := buildSay(format, args...)
	s.Render()
	Pause()
}
