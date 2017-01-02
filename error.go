package gadu

/*
#include <errno.h>
#include <string.h>
*/
import "C"

const (
	// EACCES means Access Denied
	EACCES = C.EACCES
	// EFAULT means Fault
	EFAULT = C.EFAULT
)

// GGError is a custom struct representing a GG error
type GGError struct {
	errno int
}

// NewGGError creates new GG error from errno
func NewGGError(errno int) error {
	return &GGError{errno}
}

// Fault for critical failures
var Fault = NewGGError(EFAULT)

// AccessDeniedError when unable to login etc.
var AccessDeniedError = NewGGError(EACCES)

func (e *GGError) Error() string {
	return C.GoString(C.strerror((C.int)(e.errno)))
}
