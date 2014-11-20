package mcclintock

import (
	//"bytes"
	//"fmt"
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
	style chalk.Style
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
	row := &row{}
	for _, rValue := range rowStrings {
		row.addCell(cell{rValue, len(rValue)})
	}
	for len(row) < t.columns {
		row.addCell(cell{})
	}
	t.rows = append(t.rows, row)
	return nil
}
