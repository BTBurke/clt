package clt

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// Justification sets the default placement of text inside each cell of a column
type Justification int

// Column justification flags
const (
	Left Justification = iota
	Center
	Right
)

// Cell represents a cell in the table.  Most often you'll create a cell using StyledCell
// in conjuction with AddStyledRow
type Cell struct {
	value string
	width int
	style *Style
}

// Title is a special cell that is rendered at the center top of the table that can contain
// its own styling.
type Title struct {
	value string
	width int
	style *Style
}

// Row is a row of cells in a table.  You want to use AddRow or AddStyledRow to create one.
type Row struct {
	cells []Cell
}

// Col is a column of a table.  Use ColumnHeaders, ColumnStyles, etc. to adjust default
// styling and header properties.  You can always override a particular cell in a column
// by passing in a different Cell style when you AddStyledRow.
type Col struct {
	index         int
	naturalWidth  int
	computedWidth int
	wrap          bool
	style         *Style
	justify       Justification
}

// Table is a table output to the console.  Use NewTable to construct the table with sensible defaults.
// Tables detect the terminal width and step through a number of rendering strategies to intelligently
// wrap column information to fit within the available space.
type Table struct {
	title     Title
	columns   []Col
	headers   []Cell
	rows      []Row
	pad       int
	maxWidth  int
	maxHeight int

	writer io.Writer
}

// TableOption is a function that sets an option on a table
type TableOption func(t *Table) error

// MaxHeight sets the table maximum height that can be used for pagination of
// long tables
func MaxHeight(h int) TableOption {
	return func(t *Table) error {
		t.maxHeight = h
		return nil
	}
}

// MaxWidth sets the max width of the table.  The actual max width will be set to the
// smaller of this number of the detected width of the terminal.  Very small max widths can
// be a problem because the layout engine will not be able to find a strategy to render the
// table.
func MaxWidth(w int) TableOption {
	return func(t *Table) error {
		if t.maxWidth > w {
			t.maxWidth = w
		}
		return nil
	}
}

// Magic from the go source for ssh/terminal to find terminal size.  Because it is
// difficult to test this implementation across os/terminal combinations, you can skip
// this check by setting Table.SkipTermSize=true and set the Table.maxWidth and
// Table.MaxHeight manually.  Or, Table will set appropriate conservative defaults.
func getTerminalSize() (width, height int, err error) {
	fd := os.Stdout.Fd()
	var dimensions [4]uint16

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		return -1, -1, err
	}
	return int(dimensions[1]), int(dimensions[0]), nil
}

func (r *Row) addCell(c Cell) {
	r.cells = append(r.cells, c)
}

// AddRow adds a new row to the table given an array of strings for each column's
// content.  You can set styles on a row by using AddStyledRow instead.  If you add more cells
// than available columns, the cells will be silently truncated.  If there are fewer values than columns,
// the remaining columns will be empty.
func (t *Table) AddRow(rowStrings ...string) *Table {
	newRow := Row{}
	for i, rValue := range rowStrings {
		if i >= len(t.columns) {
			break
		}
		newRow.addCell(Cell{value: rValue, width: len(rValue), style: t.columns[i].style})
	}
	for len(newRow.cells) < len(t.columns) {
		newRow.addCell(Cell{value: "", width: 0, style: Styled(Default)})
	}
	t.rows = append(t.rows, newRow)
	return t
}

// AddStyledRow adds a new row to the table with custom styles for each Cell.  If you add more cells
// than available columns, the cells will be silently truncated.  If there are fewer values than columns,
// the remaining columns will be empty.
func (t *Table) AddStyledRow(cells ...Cell) *Table {
	newRow := Row{}
	for i, cell1 := range cells {
		if i >= len(t.columns) {
			break
		}
		newRow.addCell(cell1)
	}
	for len(newRow.cells) < len(t.columns) {
		newRow.addCell(Cell{value: "", width: 0, style: Styled(Default)})
	}
	t.rows = append(t.rows, newRow)
	return t
}

// StyledCell returns a new cell with a custom style for use with AddStyledRow
func StyledCell(v string, sty *Style) Cell {
	return Cell{value: v, width: len(v), style: sty}
}

// ColumnStyles sets the default styles for each column in the row except
// the column headers.
func (t *Table) ColumnStyles(styles ...*Style) *Table {
	for i, sty := range styles {
		if i >= len(t.columns) {
			return t
		}
		t.columns[i].style = sty
	}
	return t
}

// Title sets the title for the table.  The default style is bold, but can
// be changed by passing your own styles
func (t *Table) Title(s string, styles ...Styler) *Table {
	var sty *Style
	switch {
	case len(styles) > 0:
		sty = Styled(styles...)
	default:
		sty = Styled(Bold)
	}
	t.title = Title{value: s, width: len(s), style: sty}
	return t
}

// ColumnHeaders sets the column headers with an array of strings
// The default style is Underline and Bold.  This can be changed through
// a call to SetColumnHeaderStyles.
func (t *Table) ColumnHeaders(headers ...string) *Table {
	for i, header := range headers {
		if i >= len(t.columns) {
			return t
		}
		t.headers[i].value = header
		t.headers[i].style = Styled(Bold, Underline)
		t.headers[i].width = len(header)
	}
	return t
}

// ColumnHeaderStyles sets the column header styles.
func (t *Table) ColumnHeaderStyles(styles ...*Style) *Table {
	for i, style := range styles {
		if i > len(t.columns) {
			return t
		}
		t.headers[i].style = style
	}
	return t
}

// Justification sets the justification of each column.  If you pass more justifications
// than the number of columns they will be silently dropped.
func (t *Table) Justification(cellJustifications ...Justification) *Table {
	for i, just := range cellJustifications {
		if i > len(t.columns) {
			return t
		}
		t.columns[i].justify = just
	}
	return t
}

// NewTable creates a new table with a given number of columns, setting the default
// justfication to left, and attempting to detect the existing terminal size to
// set size defaults.
func NewTable(numColumns int, options ...TableOption) *Table {
	w, h, err := getTerminalSize()
	if err != nil || w == 0 || h == 0 {
		w = 80
		h = 25
	}

	// Fill with defaults to skip complicated bounds checking on
	// changing justify or row styles
	defaultColumns := make([]Col, numColumns)
	emptyHeaders := make([]Cell, numColumns)
	for i := 0; i < numColumns; i++ {
		defaultColumns[i].index = i
		defaultColumns[i].style = Styled(Default)
		defaultColumns[i].justify = Left
		defaultColumns[i].wrap = false
	}

	t := &Table{
		columns:   defaultColumns,
		maxWidth:  w,
		maxHeight: h,
		headers:   emptyHeaders,
		pad:       1,
		title:     Title{value: "", width: 0, style: Styled(Default)},
		writer:    os.Stdout,
	}
	for _, opt := range options {
		opt(t)
	}
	return t
}

// Show will render the table using the headers, title, and styles previously
// set.
func (t *Table) Show() {
	tableAsString := t.AsString()
	fmt.Fprintf(t.writer, tableAsString)

}

// SetWriter sets the output writer if not writing to Stdout
func (t *Table) SetWriter(w io.Writer) {
	t.writer = w
}

// AsString returns the rendered table as a string instead of immediately writing to the configured writer
func (t *Table) AsString() string {
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
func renderHeaders(cells []Cell, cols []Col, pad int) string {
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
func renderRow(cells []Cell, cols []Col, pad int) string {
	wrappedLinesCount := make([]int, len(cells))

	for i, cell1 := range cells {
		wrappedL := wrap(cell1.value, cols[i].computedWidth)
		wrappedLinesCount[i] = len(wrappedL)
	}
	_, totalLines := max(wrappedLinesCount)
	lines := make([]bytes.Buffer, totalLines)

	for cellN, cellV := range cells {
		// override column style with cell style if different
		var sty *Style
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
func renderCell(s string, width int, pad int, sty *Style, justify Justification) string {
	switch justify {
	case Left:
		return justLeft(s, width, pad, sty)
	case Center:
		return justCenter(s, width, pad, sty)
	case Right:
		return justRight(s, width, pad, sty)
	}
	return ""
}

// justCenter is center-justified text with padding and style
func justCenter(s string, width int, pad int, sty *Style) string {
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
func justLeft(s string, width int, pad int, sty *Style) string {
	contentLen := len(s)
	onRight := width - contentLen
	if onRight < 0 {
		onRight = 0
	}
	return fmt.Sprintf("%s%s%s", spaces(pad), sty.ApplyTo(s), spaces(onRight+pad))
}

// justRight is right-justified text with padding and style
func justRight(s string, width int, pad int, sty *Style) string {
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
	return fmt.Errorf("no table rendering strategy suitable")
}

// simpleStrategy sets all column widths to their natural width.
// Successful if the whole table fits inside maxWidth (including pad)
func simpleStrategy(t *Table) bool {
	natWidths := extractNatWidth(t)
	colWPadded := mapAdd(natWidths, 2*t.pad)
	totalWidth := sum(colWPadded)

	if totalWidth <= t.maxWidth {
		for i := range t.columns {
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
	tableMaxW := t.maxWidth - 2*len(t.columns)*t.pad
	wrapW := tableMaxW - sumWithoutIndex(naturalWidths, maxI)
	if wrappedWidthOk(wrapW, maxW) {
		for i := range t.columns {
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
// from []Col
func extractNatWidth(t *Table) []int {
	out := make([]int, len(t.columns))
	for i, col := range t.columns {
		out[i] = col.naturalWidth
	}
	return out
}

// convenience function for extracting computed width as []int
// from []Col
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
