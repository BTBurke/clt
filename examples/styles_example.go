package main

import (
	"fmt"
	"github.com/BTBurke/mcclintock"
)

func main() {
	r := mcclintock.Text().Red("RED")
	fmt.Println(r)
	fmt.Printf("\x1b[31mThis sentence is red\x1b[39m")
	
}