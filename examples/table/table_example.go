package main

import (
	"fmt"
	"strings"

	"github.com/BTBurke/clt"
)

func main() {
	sepLine := strings.Repeat("=", 15)

	// A simple table that should fit in a standard terminal width 80
	fmt.Printf("\n\n\n%s Simple Table Example %s\n\n\n", sepLine, sepLine)
	SimpleTable()

	// A table with long content that needs to be wrapped.  The table
	// library has several strategies for fitting the content into the
	// available terminal space.
	fmt.Printf("\n\n\n%s Wrapped Table Example %s\n\n\n", sepLine, sepLine)
	WrappedTable()

	// Tables can be styled many ways, using the clt.Styled library
	fmt.Printf("\n\n\n%s Styled Table Example %s\n\n\n", sepLine, sepLine)
	StyledTable()

	fmt.Printf("\n\n")
}

func SimpleTable() {
	// A 3-column table
	t := clt.NewTable(3)

	// Set the headers and title
	t.ColumnHeaders("Column1", "Column2", "Column3")
	t.Title("Simple Example Table")

	// Add some rows
	t.AddRow("Col1 Line1", "Col2 Line1", "Col3 Line1")
	t.AddRow("Col1 Line2", "Col2 Line2", "Col3 Line2")

	// Print the table
	t.Show()
}

func WrappedTable() {
	// A 3-column table
	// Force the terminal size to be small to see wrapping behavior.
	// Normally, the terminal width is detected automatically, but
	// you can set the table MaxWidth explicitly when desired.
	t := clt.NewTable(3, clt.MaxWidth(50), clt.Spacing(2)).
		ColumnHeaders("Column1", "Column2", "Column3").
		Title("Wrapped Example Table")

	// Add some rows
	t.AddRow("Col1 Line1", "Col2 Line1", "This is a pretty long description.")
	t.AddRow("Col1 Line2", "Col2 Line2", "This is another longish one.")

	// Print the table
	t.Show()
}

func StyledTable() {
	// A 3-column table
	t := clt.NewTable(2)

	// Set the headers and title
	t.ColumnHeaders("Status", "Reason")
	t.Title("Styled Example Table")

	// Set styles for each column
	t.ColumnStyles(clt.Styled(clt.Green), clt.Styled(clt.Default))

	// Add some rows.  The OK will be green.
	t.AddRow("OK", "Everything worked")

	// Add another row with custom styling to override the green column
	t.AddStyledRow(clt.StyledCell("FAIL", clt.Styled(clt.Red)), clt.StyledCell("Something bad happened", clt.Styled(clt.Default)))

	// Print the table
	t.Show()
}
