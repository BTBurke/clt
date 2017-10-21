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
	fmt.Printf("This is %s text\n", clt.Styled(clt.Blue, clt.Underline).ApplyTo("blue and underlined"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Blue, clt.Background(clt.White)).ApplyTo("blue on a white background"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Italic).ApplyTo("italic"))
	fmt.Printf("This is %s text\n", clt.Styled(clt.Bold).ApplyTo("bold"))
}
```
![console output](https://s3.amazonaws.com/btburke-github/styles_example.png)

The general operation of the style function is to first call `clt.Styled(<Style1>, <Style2>, ...)`.  This creates a style that can then be applied to a string via the `.ApplyTo(<string>)` method.

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
	t.SetTitle("Hockey Standings")

	// Set column headers
	t.SetColumnHeaders("Team", "Points", "W", "L", "OT")

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

CLT provides two kinds of progress bars:

* *Spinner* - Useful for when you want to show progress but don't know exactly when an action will complete

* *Bar* - Useful when you have a defined number of iterations to completion and you can update progress during processing

#### Example:  

See [examples/progress_example.go](https://github.com/BTBurke/clt/blob/master/examples/progress_example.go) for the example in the screencast below.

![console output](https://s3.amazonaws.com/btburke-github/progress-ex.gif)

Progress bars use go routines to update the progress status while your app does other processing.  Remember to close out the progress element with either a call to `Success()` or `Fail()` to terminate this routine.




