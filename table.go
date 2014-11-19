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
