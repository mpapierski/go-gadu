package gadu

/*
#include <errno.h>
#include <string.h>
*/
import "C"

const (
	EACCES = C.EACCES
	EFAULT = C.EFAULT
)

type GGError struct {
	errno int
}

func NewGGError(errno int) error {
	return &GGError{errno}
}

var Fault = NewGGError(EFAULT)
var AccessDeniedError = NewGGError(EACCES)

func (e *GGError) Error() string {
	return C.GoString(C.strerror((C.int)(e.errno)))
}
