package mcclintock

import (
	//"bytes"
	"fmt"
	"github.com/ttacon/chalk"
	"os"
	//"os/exec"
	//"strconv"
	//"strings"
	"syscall"
	"unsafe"
)

type cell struct {
	value string
	width int
	style chalk.Style
}

type title struct {
	value string
	width int
	style chalk.TextStyle
}

type row struct {
	cells []cell
}

type table struct {
	title        title
	columns      int
	headers      []cell
	rows         []row
	ShowLines    bool
	Indent       int
	MaxWidth     int
	MaxHeight    int
	justify      []string
	SkipTermSize bool
	tableStyle   []chalk.Style
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
func (t *table) AddRow(rowStrings []string) error {
	if len(rowStrings) > t.columns {
		return fmt.Errorf("Received %v columns but table only has %v columns.", len(rowStrings), t.columns)
	}
	newRow := row{}
	for _, rValue := range rowStrings {
		newRow.addCell(cell{value: rValue, width: len(rValue)})
	}
	for len(newRow.cells) < t.columns {
		newRow.addCell(cell{})
	}
	t.rows = append(t.rows, newRow)
	return nil
}

// SetTableStyle sets the default styles for each column in the row except
// the column headers.
func (t *table) SetTableStyle(styles ...chalk.Style) error {
	if len(styles) > t.columns {
		return fmt.Errorf("Received %v column styles but table only has %v columns.", len(styles), t.columns)
	}

}

// SetTitle sets the title for the table.  The default style is bold, but can
// be changed using SetTitleStyle.
func (t *table) SetTitle(s string) {
	t.title = title{value: s, width: len(s), style: chalk.Bold}
}

// SetTitleStyle sets the font style for the title.  The default is chalk.Bold
// but can be set to any valid value of chalk.TextStyle
func (t *table) SetTitleStyle(sty chalk.TextStyle) {
	t.title.style = sty
}

// SetColumnHeaders sets the column headers

// NewTable creates a new table with a given number of columns, setting the default
// justfication to left, and attempting to detect the existing terminal size to
// set size defaults.  You can change these defaults using SetJustification and
// MaxWidth and MaxHeight properties.
func NewTable(numColumns int) *table {
	w, h, err := getTerminalSize()
	if err != nil {
		w = 80
		h = 25
	}
	just := make([]string, numColumns)
	for i := 0; i < numColumns; i++ {
		just[i] = "l"
	}
	return &table{columns: numColumns, justify: just, MaxWidth: w, MaxHeight: h}
}
