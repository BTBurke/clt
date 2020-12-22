// +build windows

package clt

import (
	"log"
	"syscall"
)

func init() {
	enableANSI()
}

var (
	dll            = syscall.MustLoadDLL("kernel32")
	setConsoleMode = dll.MustFindProc("SetConsoleMode")
)

func setInputConsoleMode(h syscall.Handle, m uint32) error {
	r, _, err := setConsoleMode.Call(uintptr(h), uintptr(m))
	if r == 0 {
		return err
	}
	return nil
}

func enableANSI() {
	h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		log.Printf("error getting stdout handle")
	}
	if err := setInputConsoleMode(h, 4); err != nil {
		log.Printf("error setting ANSI handling: %s", err)
	}
}
