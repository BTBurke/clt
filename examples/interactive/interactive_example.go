package main

import (
	"fmt"

	"github.com/BTBurke/clt"
)

func main() {
	fmt.Println("------- Pagination Example --------")
	Pagination()

	fmt.Println("\n\n\n\n\n------- Yes/No Example --------")
	AskYesNo()

	fmt.Println("\n\n\n\n\n------- Choice Example --------")
	ChooseFromOptions()

	fmt.Println("\n\n\n\n\n------- Password Example -------")
	PasswordPrompt()
}

func Pagination() {
	i := clt.NewInteractiveSession()
	i.Say("This can be a really long screed that needs to be paginated.  Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	i.Pause()
	i.Say("And I could continue with more stuff...")
}

func AskYesNo() {
	i := clt.NewInteractiveSession()
	resp := i.Say("This is an example of asking a yes/no question and maybe you want to add a warning as well.").
		Warn("%s", clt.SStyled("Bad things can happen if you do this!", clt.Bold)).
		AskYesNo("Do you really want to do this?", "n")
	switch {
	case clt.IsYes(resp):
		i.Say("Ok, I'll go ahead and do that.")
	case clt.IsNo(resp):
		i.Say("Good idea, let's do that later.")
	}
}

func ChooseFromOptions() {
	choices := map[string]string{
		"a":     "Do task a",
		"b":     "Do task b",
		"abort": "Let's get out of here",
	}
	i := clt.NewInteractiveSession()
	resp := i.Say("You can also create a list of options and let them select from the list.").
		AskFromTable("Pick a choice from the table", choices, "a")
	i.Say("Ok, let's do: %s", choices[resp])
}

func PasswordPrompt() {
	i := clt.NewInteractiveSession()
	pw := i.AskPassword()
	i.Say("Shhh! Don't tell anyone your password is %s", pw)
}
