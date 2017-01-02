// +build !darwin

package gadu

import "syscall"

func syscallSelect(nfd int, r *syscall.FdSet, w *syscall.FdSet, e *syscall.FdSet, timeout *syscall.Timeval) (n int, err error) {
	return syscall.Select(nfd, r, w, e, timeout)
}
