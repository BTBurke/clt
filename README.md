[![](https://godoc.org/github.com/BTBurke/clt?status.svg)](http://godoc.org/github.com/BTBurke/clt)

CLT - Command Line Tools for Go
===
CLT is a toolset for building elegant command line interfaces in Go.  CLT includes elements like styled text, tables, selections from a list, and more so that you can quickly build CLIs with interactive elements without the hassle of dealing with formatting all these yourself.

Go Doc documentation is available at [Godoc] and examples are located in the examples directory.  This readme strives to show you the major features.

## Styled Text

```go
package main

import (
	"fmt"
	"github.com/BTBurke/clt"
)

func main() {
	fmt.Printf("This is %s text\n", clt.Styled(clt.Red).ApplyTo("red"))
	fmt.Printf("This is %s text\n", clt.SStyled("blue and underlined", clt.Blue, clt.Underline))
	fmt.Printf("This is %s text\n", clt.SStyled("blue on a white background", clt.Blue, clt.Background(clt.White))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Italic).ApplyTo("italic"))
	fmt.Printf("This is %s text\n", clt.SStyled("bold", clt.Bold)
}
```
![console output](https://s3.amazonaws.com/btburke-github/styles_example.png)

The general operation of the style function is to first call `clt.Styled(<Style1>, <Style2>, ...)`.  This creates a style that can then be applied to a string via the `.ApplyTo(<string>)` method.  A shortcut method `clt.SStyled("string", styles...)` can help eliminate some of the boilerplate.

## Tables

CLT provides an easy-to-use library for building text tables.  It provides layout algorithms for multi-column tables and the ability to style each column or individual cells using clt.Styled.

Tables detect the terminal width and intelligently decide how cell contents should be wrapped to fit on screen.
```go
package main

import "github.com/BTBurke/clt"

func main() {

	// Create a table with 3 columns
	t := clt.NewTable(5)

	// Add a title
	t.Title("Hockey Standings")

	// Set column headers
	t.ColumnHeaders("Team", "Points", "W", "L", "OT")

	// Add some rows
	t.AddRow("Washington Capitals", "42", "18", "11", "6")
	t.AddRow("NJ Devils", "31", "12", "18", "7")

	// Render the table
	t.Show()
}
```

Produces:

![console output](https://s3.amazonaws.com/btburke-github/simple-table.png)

#### More examples
See [examples/table_example.go](https://github.com/BTBurke/clt/blob/master/examples/table_example.go) for more examples.  Also, see the GoDoc for the details of the table library.

## Progress Bars

CLT provides three kinds of progress indicators:

* *Spinner* - Useful for when you want to show progress but don't know exactly when an action will complete

* *Bar* - Useful when you have a defined number of iterations to completion and you can update progress during processing

* *Loading* - Useful for when you are making a remote call and you want to give a visual indication that something is going on in the background, but you want it to disappear as soon as the call ends.  It also has a configurable delay so that the loading indicator will only appear when the call takes longer than the delay to complete.

#### Example:  

See [examples/progress_example.go](https://github.com/BTBurke/clt/blob/master/examples/progress_example.go) for the example in the screencast below.

![console output](https://s3.amazonaws.com/btburke-github/progress-ex-20171025.gif)

Progress bars use go routines to update the progress status while your app does other processing.  Remember to close out the progress element with either a call to `Success()` or `Fail()` to terminate this routine.

## Interactive Sessions

CLT provides building blocks to create interactive sessions, giving you flexible functions to ask the user for input.

See [examples/interactive_example.go](https://github.com/BTBurke/clt/blob/master/examples/interactive_example.go) for examples of creating interactive interfaces.

#### Interactions

| Interaction | Use |
| ------- | ----------- |
| Ask | Ask for a response with optional validation |
| AskWithDefault | Ask with a preconfigured default value |
| AskWithHint | Ask with a hint that shows how the input should be formatted |
| AskPassword | Ask for a password without any echo to the terminal while the user types |
| AskYesNo | Ask a yes or no question with a default to either |
| AskFromTable | User picks an option from a table of possibilities |
| Pause | Paginate some output with `Press [Enter] to continue.` |
| PauseWithPrompt | Paginate some output with a custom prompt to continue |
| Warn | Issue a warning that is visually separated from other text by style changes |
| Error | Show an error message and exit the process |
| Say | Thin wrapper around `fmt.Printf` that helps build interactive sessions in a fluent style and take care of common spacing issues |



