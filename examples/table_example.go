package main

import (
	"fmt"
	"github.com/BTBurke/clt"
)

func main() {

	// A simple table that should fit in a standard terminal width 80
	fmt.Printf("\n")
	SimpleTable()

	// A table with long content that needs to be wrapped.  The table
	// library has several strategies for fitting the content into the
	// available terminal space.
	fmt.Printf("\n\n")
	WrappedTable()

	// Tables can be styled many ways, using the clt.Style library
	fmt.Printf("\n\n")
	StyledTable()
}

func SimpleTable() {
	// A 3-column table
	t := clt.NewTable(3)

	// Set the headers and title
	t.SetColumnHeaders("Column1", "Column2", "Column3")
	t.SetTitle("Simple Example Table")

	// Add some rows
	t.AddRow("Col1 Line1", "Col2 Line1", "Col3 Line1")
	t.AddRow("Col1 Line2", "Col2 Line2", "Col3 Line2")

	// Print the table
	t.Show()
}

func WrappedTable() {
	// A 3-column table
	t := clt.NewTable(3)

	// Set the headers and title
	t.SetColumnHeaders("Column1", "Column2", "Column3")
	t.SetTitle("Wrapped Example Table")

	// Add some rows
	t.AddRow("Col1 Line1", "Col2 Line1", "This is a pretty long description.")
	t.AddRow("Col1 Line2", "Col2 Line2", "This is another longish one.")

	// Force the terminal size to be small to see wrapping behavior.
	// Normally, the terminal width is detected automatically, but
	// you can set the table MaxWidth explicitly when desired.
	t.MaxWidth = 50

	// Print the table
	t.Show()
}

func StyledTable() {
	// A 3-column table
	t := clt.NewTable(2)

	// Set the headers and title
	t.SetColumnHeaders("Status", "Reason")
	t.SetTitle("Styled Example Table")

	// Set styles for each column
	t.SetColumnStyles(clt.Style(clt.Green), clt.Style(clt.Default))

	// Add some rows.  The OK will be green.
	t.AddRow("OK", "Everything worked")

	// Add another row with custom styling to override the green column
	t.AddStyledRow(clt.Cell("FAIL", clt.Style(clt.Red)), clt.Cell("Something bad happened", clt.Style(clt.Default)))

	// Print the table
	t.Show()
}
