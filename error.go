package gadu

/*
#include <errno.h>
#include <string.h>
*/
import "C"

type GGError struct {
	errno int
}

func NewGGError(errno int) error {
	return &GGError{errno}
}

func (e *GGError) Error() string {
	return C.GoString(C.strerror((C.int)(e.errno)))
}
