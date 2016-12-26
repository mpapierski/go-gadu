package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>
*/
import "C"

func Version() string {

	return C.GoString(C.gg_libgadu_version())
}
