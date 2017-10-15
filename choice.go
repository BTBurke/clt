package clt

import (
	"fmt"
	"strconv"
)

func Choice(prompt string) string {
	c := NewConsoleInput()
	c.Prompt = prompt
	return c.Get()
}

func ChoiceWithDefault(prompt string, def string) string {
	c := NewConsoleInput()
	c.Prompt = prompt
	c.Default = def
	return c.Get()
}

func ChoiceWithValidation(prompt string, hint string, v ValidationFunc) string {
	c := NewConsoleInput()
	c.Prompt = prompt
	c.ValHint = hint
	c.AskValidate(v)
	return c.Response
}

func ChoiceFromTable(prompt string, choices []string, def int) string {

	t := NewTable(2)
	for key, choice := range choices {
		t.AddRow(strconv.Itoa(key), choice)
	}
	tAsString := t.renderTableAsString()

	c := NewConsoleInput()
	c.Prompt = fmt.Sprintf("%s\n%s\nChoice [%s]: ", prompt, tAsString, choices[def])
	c.AskValidate(OptionValidator(choices))
	return c.Response
}
