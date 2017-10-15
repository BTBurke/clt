package clt

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// Column justification flags
const (
	jLeft int = iota
	jCenter
	jRight
)

type cell struct {
	value string
	width int
	style *style
}

type title struct {
	value string
	width int
	style *style
}

type row struct {
	cells []cell
}

type col struct {
	index         int
	naturalWidth  int
	computedWidth int
	wrap          bool
	style         *style
	justify       int
}

// Table is a table output to the console.  Use NewTable to construct the table with sensible defaults.
// Tables detect the terminal width and step through a number of rendering strategies to intelligently
// wrap column information to fit within the available space.
type Table struct {
	title        title
	columns      []col
	headers      []cell
	rows         []row
	pad          int
	MaxWidth     int
	MaxHeight    int
	SkipTermSize bool
}

// Magic from the go source for ssh/terminal to find terminal size.  Because it is
// difficult to test this implementation across os/terminal combinations, you can skip
// this check by setting Table.SkipTermSize=true and set the Table.MaxWidth and
// Table.MaxHeight manually.  Or, Table will set appropriate conservative defaults.
func getTerminalSize() (width, height int, err error) {
	fd := os.Stdout.Fd()
	var dimensions [4]uint16

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		return -1, -1, err
	}
	return int(dimensions[1]), int(dimensions[0]), nil
}

func (r *row) addCell(c cell) {
	r.cells = append(r.cells, c)
}

// AddRow adds a new row to the table given an array of strings for each column's
// content.  You can set styles on this particular row by a subsequent call to
// AddRowStyle.
func (t *Table) AddRow(rowStrings ...string) error {
	if len(rowStrings) > len(t.columns) {
		return fmt.Errorf("Received %v columns but table only has %v columns.", len(rowStrings), len(t.columns))
	}
	newRow := row{}
	for i, rValue := range rowStrings {
		newRow.addCell(cell{value: rValue, width: len(rValue), style: t.columns[i].style})
	}
	for len(newRow.cells) < len(t.columns) {
		newRow.addCell(cell{value: "", width: 0, style: Style(Default)})
	}
	t.rows = append(t.rows, newRow)
	return nil
}

// AddStyledRow adds a new row to the table with custom styles for each cell.
func (t *Table) AddStyledRow(cells ...cell) error {
	if len(cells) > len(t.columns) {
		return fmt.Errorf("Received %v columns but table only has %v columns.", len(cells), len(t.columns))
	}
	newRow := row{}
	for _, cell1 := range cells {
		newRow.addCell(cell1)
	}
	for len(newRow.cells) < len(t.columns) {
		newRow.addCell(cell{value: "", width: 0, style: Style(Default)})
	}
	t.rows = append(t.rows, newRow)
	return nil
}

// Cell returns a new cell with a custom style for use with AddStyledRow
func Cell(v string, sty *style) cell {
	return cell{value: v, width: len(v), style: sty}
}

// SetColumnStyles sets the default styles for each column in the row except
// the column headers.
func (t *Table) SetColumnStyles(styles ...*style) error {
	if len(styles) > len(t.columns) {
		return fmt.Errorf("Received %v column styles but table only has %v columns.", len(styles), len(t.columns))
	}
	for i, sty := range styles {
		t.columns[i].style = sty
	}
	return nil
}

// SetTitle sets the title for the table.  The default style is bold, but can
// be changed using SetTitleStyle.
func (t *Table) SetTitle(s string) {
	t.title = title{value: s, width: len(s), style: Style(Bold)}
}

// SetTitleStyle sets the font style for the title.  The default is chalk.Bold
// but can be set to any valid value of chalk.TextStyle
func (t *Table) SetTitleStyle(sty *style) {
	t.title.style = sty
}

// SetColumnHeaders sets the column headers with an array of strings
// The default style is Underline and Bold.  This can be changed through
// a call to SetColumnHeaderStyles.
func (t *Table) SetColumnHeaders(headers ...string) error {
	if len(headers) > len(t.columns) {
		return fmt.Errorf("More column headers than columns.")
	}
	for i, header := range headers {
		t.headers[i].value = header
		t.headers[i].style = Style(Bold, Underline)
		t.headers[i].width = len(header)
	}
	return nil
}

// SetColumnHeaderStyles sets the column header styles. Returns an error
// if there are more styles than the number of columns.
func (t *Table) SetColumnHeaderStyles(styles ...*style) error {
	if len(styles) > len(t.columns) {
		return fmt.Errorf("Got more styles than number of columns")
	}
	for i, style := range styles {
		t.headers[i].style = style
	}
	return nil
}

// NewTable creates a new table with a given number of columns, setting the default
// justfication to left, and attempting to detect the existing terminal size to
// set size defaults.  You can change these defaults using SetJustification and
// MaxWidth and MaxHeight properties.
func NewTable(numColumns int) *Table {
	w, h, err := getTerminalSize()
	if err != nil || w == 0 || h == 0 {
		w = 80
		h = 25
	}

	// Fill with defaults to skip complicated bounds checking on
	// changing justify or row styles
	defaultColumns := make([]col, numColumns)
	emptyHeaders := make([]cell, numColumns)
	for i := 0; i < numColumns; i++ {
		defaultColumns[i].index = i
		defaultColumns[i].style = Style(Default)
		defaultColumns[i].justify = jLeft
		defaultColumns[i].wrap = false
	}

	return &Table{
		columns:   defaultColumns,
		MaxWidth:  w,
		MaxHeight: h,
		headers:   emptyHeaders,
		pad:       1,
		title:     title{value: "", width: 0, style: Style(Default)},
	}

}

// Show will render the table using the headers, title, and styles previously
// set.
func (t *Table) Show() {
	tableAsString := t.renderTableAsString()
	fmt.Printf(tableAsString)

}

// returns the rendered table as a string
func (t *Table) renderTableAsString() string {
	err := t.computeColWidths()
	if err != nil {
		// this error should never happen with fallback overflow strategy
		log.Fatal(err)
	}
	var renderedT bytes.Buffer
	renderedT.WriteString(renderTitle(t) + "\n\n")
	renderedT.WriteString(renderHeaders(t.headers, t.columns, t.pad))
	for _, row := range t.rows {
		renderedT.WriteString(renderRow(row.cells, t.columns, t.pad))
	}
	return renderedT.String()
}

// renderTitle returns the title as a formatted string
func renderTitle(t *Table) string {
	return justCenter(t.title.value, t.width(), 0, t.title.style)
}

// renders the headers as a string
func renderHeaders(cells []cell, cols []col, pad int) string {
	wrappedLinesCount := make([]int, len(cells))

	for i, cell1 := range cells {
		wrappedL := wrap(cell1.value, cols[i].computedWidth)
		wrappedLinesCount[i] = len(wrappedL)
	}
	_, totalLines := max(wrappedLinesCount)
	lines := make([]bytes.Buffer, totalLines)

	for cellN, cellV := range cells {
		wL := wrap(cellV.value, cols[cellN].computedWidth)
		for i := 0; i < totalLines; i++ {
			switch {
			case i < len(wL):
				lines[i].WriteString(renderCell(wL[i], cols[cellN].computedWidth, pad, cellV.style, cols[cellN].justify))
			default:
				lines[i].WriteString(renderCell("", cols[cellN].computedWidth, pad, cellV.style, cols[cellN].justify))
			}
		}
	}
	var out bytes.Buffer
	for _, line := range lines {
		out.Write(line.Bytes())
		out.WriteString("\n")
	}
	return out.String()
}

// renderRow renders the row as a styled string and implements the
// wrapping of long strings where necessary
func renderRow(cells []cell, cols []col, pad int) string {
	wrappedLinesCount := make([]int, len(cells))

	for i, cell1 := range cells {
		wrappedL := wrap(cell1.value, cols[i].computedWidth)
		wrappedLinesCount[i] = len(wrappedL)
	}
	_, totalLines := max(wrappedLinesCount)
	lines := make([]bytes.Buffer, totalLines)

	for cellN, cellV := range cells {
		// override column style with cell style if different
		var sty *style
		switch {
		case cellV.style != cols[cellN].style:
			sty = cellV.style
		default:
			sty = cols[cellN].style
		}

		wL := wrap(cellV.value, cols[cellN].computedWidth)
		for i := 0; i < totalLines; i++ {
			switch {
			case i < len(wL):
				lines[i].WriteString(renderCell(wL[i], cols[cellN].computedWidth, pad, sty, cols[cellN].justify))
			default:
				lines[i].WriteString(renderCell("", cols[cellN].computedWidth, pad, sty, cols[cellN].justify))
			}
		}
	}
	var out bytes.Buffer
	for _, line := range lines {
		out.Write(line.Bytes())
		out.WriteString("\n")
	}
	return out.String()
}

// renderCell renders the cell as a string using the correct justification
func renderCell(s string, width int, pad int, sty *style, justify int) string {
	switch justify {
	case jLeft:
		return justLeft(s, width, pad, sty)
	case jCenter:
		return justCenter(s, width, pad, sty)
	case jRight:
		return justRight(s, width, pad, sty)
	}
	return ""
}

// justCenter is center-justified text with padding and style
func justCenter(s string, width int, pad int, sty *style) string {
	contentLen := len(s)
	onLeft := (width - contentLen) / 2
	if onLeft < 0 {
		onLeft = 0
	}
	onRight := width - contentLen - onLeft
	if onRight < 0 {
		onRight = 0
	}
	return fmt.Sprintf("%s%s%s", spaces(onLeft+pad), sty.ApplyTo(s), spaces(onRight+pad))
}

// justLeft is left-justified text with padding and style
func justLeft(s string, width int, pad int, sty *style) string {
	contentLen := len(s)
	onRight := width - contentLen
	if onRight < 0 {
		onRight = 0
	}
	return fmt.Sprintf("%s%s%s", spaces(pad), sty.ApplyTo(s), spaces(onRight+pad))
}

// justRight is right-justified text with padding and style
func justRight(s string, width int, pad int, sty *style) string {
	contentLen := len(s)
	onLeft := width - contentLen
	if onLeft < 0 {
		onLeft = 0
	}
	return fmt.Sprintf("%s%s%s", spaces(onLeft+pad), sty.ApplyTo(s), spaces(pad))
}

// wrap will break long lines on breakpoints space, :, ., /, \, -.  If
// line is too long without breakpoints, will do dumb wrap at width w.
func wrap(s string, w int) []string {
	var out []string
	var wrapped string
	rem := s
	for len(rem) > 0 {
		wrapped, rem = wrapSubString(rem, w, " :.-/\\")
		out = append(out, wrapped)
	}
	return out
}

// wrapSubString - don't call directly. Works with wrap to recursively
// split a string at the specified breakpoints.
func wrapSubString(s string, w int, breakpts string) (wrapped string, remainder string) {

	if len(s) <= w {
		return strings.TrimSpace(s), ""
	}

	ind := strings.LastIndexAny(s[0:w], breakpts)
	switch {
	case ind > 0:
		return strings.TrimSpace(s[0 : ind+1]), strings.TrimSpace(s[ind+1 : len(s)])
	case ind == -1:
		return strings.TrimSpace(s[0:w]), strings.TrimSpace(s[w:len(s)])
	}
	return "", ""
}

// spaces is a convenience function to get n spaces repeated
func spaces(n int) string {
	return strings.Repeat(" ", n)
}

// width returns the full table computed width including padding
func (t *Table) width() int {
	return sum(extractComputedWidth(t)) + len(t.columns)*2*t.pad
}

// automagically determine column widths.  See if it can fit inside
// max width. If not, make intelligent guess about which should be
// made multi-line
func (t *Table) computeColWidths() error {
	computeNaturalWidths(t)
	switch {
	case simpleStrategy(t):
		return nil
	case wrapWidestStrategy(t):
		return nil
	case overflowStrategy(t):
		return nil
	}
	return fmt.Errorf("No table rendering strategy suitable.")
}

// simpleStrategy sets all column widths to their natural width.
// Successful if the whole table fits inside MaxWidth (including pad)
func simpleStrategy(t *Table) bool {
	natWidths := extractNatWidth(t)
	colWPadded := mapAdd(natWidths, 2*t.pad)
	totalWidth := sum(colWPadded)

	if totalWidth <= t.MaxWidth {
		for i, _ := range t.columns {
			t.columns[i].computedWidth = natWidths[i]
		}
		return true
	}
	return false
}

// wrapWidestStrategy wraps the column with the largest natural width.
// Successful if the wrapped width >50% of natural width
func wrapWidestStrategy(t *Table) bool {
	naturalWidths := extractNatWidth(t)
	maxI, maxW := max(naturalWidths)
	tableMaxW := t.MaxWidth - 2*len(t.columns)*t.pad
	wrapW := tableMaxW - sumWithoutIndex(naturalWidths, maxI)
	if wrappedWidthOk(wrapW, maxW) {
		for i, _ := range t.columns {
			switch i {
			case maxI:
				t.columns[i].computedWidth = wrapW
				t.columns[i].wrap = true
			default:
				t.columns[i].computedWidth = t.columns[i].naturalWidth
			}
		}
		return true
	}
	return false
}

// overflowStrategy is the fallback if no other strategy makes the
// table fit within the natural width. Sets all columns to their
// natural width and lets the terminal wrap the lines.
func overflowStrategy(t *Table) bool {
	for i, col := range t.columns {
		t.columns[i].computedWidth = col.naturalWidth
	}
	return true
}

// convenience function for extracting natural width as []int
// from []col
func extractNatWidth(t *Table) []int {
	out := make([]int, len(t.columns))
	for i, col := range t.columns {
		out[i] = col.naturalWidth
	}
	return out
}

// convenience function for extracting computed width as []int
// from []col
func extractComputedWidth(t *Table) []int {
	out := make([]int, len(t.columns))
	for i, col := range t.columns {
		out[i] = col.computedWidth
	}
	return out
}

// computes natural column widths and stores in table.columns.naturalWidth
func computeNaturalWidths(t *Table) {
	maxColW := make([]int, len(t.columns))

	for _, row := range t.rows {
		for col, cell := range row.cells {
			if cell.width > maxColW[col] {
				maxColW[col] = cell.width
			}
		}
	}

	for col, header := range t.headers {
		if header.width > maxColW[col] {
			maxColW[col] = header.width
		}
	}

	for i, natWidth := range maxColW {
		t.columns[i].naturalWidth = natWidth
	}
}

func sum(n []int) int {
	total := 0
	for _, num := range n {
		total += num
	}
	return total
}

// sum of all values except value at index
func sumWithoutIndex(n []int, index int) int {
	total := 0
	for idx, num := range n {
		switch idx {
		case index:
			continue
		default:
			total += num
		}
	}
	return total
}

// wrappedWidthOk true if wrapped width >50% natural width
func wrappedWidthOk(wrapW int, naturalW int) bool {
	if float64(wrapW)/float64(naturalW) >= 0.5 {
		return true
	}
	return false
}

// positive numbers only, returns index of first logical max
func max(n []int) (index int, biggest int) {
	for i, num := range n {
		if num > biggest {
			biggest = num
			index = i
		}
	}
	return
}

// add inc to every number in n, return new array
func mapAdd(n []int, inc int) []int {
	ret := make([]int, len(n))
	for i, num := range n {
		ret[i] = num + inc
	}
	return ret
}
