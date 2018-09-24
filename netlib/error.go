package netlib

import (
	"net"
	"os"
	"syscall"
)

func unboxSysOpError(err error) syscall.Errno {
	if e2, ok := err.(*net.OpError); ok {
		if e3, ok := e2.Err.(*os.SyscallError); ok {
			if errno, ok := e3.Err.(syscall.Errno); ok {
				return errno
			}
		}
	}
	return 0
}
func ConnResetError(err error) bool {
	e := unboxSysOpError(err)
	return e == syscall.ECONNRESET || e == syscall.WSAECONNRESET
}
