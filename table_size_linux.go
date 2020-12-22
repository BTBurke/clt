// +build linux

package clt

import (
	"os"
	"syscall"
	"unsafe"
)

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
