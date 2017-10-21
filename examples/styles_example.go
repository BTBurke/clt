package main

import (
	"fmt"

	"github.com/BTBurke/clt"
)

func main() {
	fmt.Printf("This is %s text\n", clt.Styled(clt.Red).ApplyTo("red"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Blue, clt.Underline).ApplyTo("blue and underlined"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Blue, clt.Background(clt.White)).ApplyTo("blue on a white background"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Italic).ApplyTo("italic"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Bold).ApplyTo("bold"))
	fmt.Printf("This is %s text\n", clt.SStyled("red underline", clt.Red, clt.Underline))
}
