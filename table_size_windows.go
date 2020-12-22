// +build windows

package clt

import "fmt"

func getTerminalSize() (width, height int, err error) {
	return -1, -1, fmt.Errorf("fuck windows")
}
