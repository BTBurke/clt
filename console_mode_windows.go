// +build windows

package clt

import (
	"log"
	"os"
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
	log.Printf("enabled ANSI handling")
	h := syscall.Handle(os.Stdin.Fd())
	if err := setInputConsoleMode(h, 0x0004); err != nil {
		log.Printf("error setting ANSI handling: %s", err)
	}
}
