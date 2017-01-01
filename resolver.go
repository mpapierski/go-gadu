package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>

extern int startGoResolver(int* fd, void** data, char* name);
extern void cleanupGoResolver(void**, int);

static int startGoResolver_wrapper(int* fd, void** data, const char* name) {
    // Another wrapper because cgo doesnt do "const"
	return startGoResolver(fd, data, (char*)name);
}

static void SetCustomResolver() {
	gg_global_set_custom_resolver(startGoResolver_wrapper, cleanupGoResolver);
}
*/
import "C"
import "unsafe"
import "syscall"

//export startGoResolver
func startGoResolver(fd *C.int, data *unsafe.Pointer, name *C.char) C.int {
	fds, err := syscall.Socketpair(0, 0)
	*fd = rd.Fd()
	return 0
}

//export cleanupGoResolver
func cleanupGoResolver(*unsafe.Pointer, C.int) {
}

func init() {
	C.SetCustomResolver()
}
