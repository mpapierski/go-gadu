package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>
#include <sys/socket.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <unistd.h>

extern int startGoResolver(int fd, char* name);
extern void cleanupGoResolver();

struct resolve_data {
	int pipes[2];
};

static void cleanupGoResolver_wrapper(void** priv_data, int force) {
	struct resolve_data* data = (struct resolve_data*)*priv_data;
	close(data->pipes[0]);
	free(data);
}

static int startGoResolver_wrapper(int* fd, void** data, const char* name) {
	struct resolve_data* priv = (struct resolve_data*)malloc(sizeof(struct resolve_data));
	if (socketpair(AF_LOCAL, SOCK_STREAM, 0, priv->pipes) == -1) {
		return -1;
	}
	*fd = priv->pipes[0];
	*data = priv;
    // Another wrapper because cgo doesnt do "const"
	return startGoResolver(priv->pipes[1], (char*)name);
}

static void SetCustomResolver() {
	gg_global_set_custom_resolver(startGoResolver_wrapper, cleanupGoResolver_wrapper);
}
*/
import "C"
import (
	"log"
	"net"
	"syscall"
	"unsafe"
)

//export startGoResolver
func startGoResolver(fd C.int, name *C.char) C.int {
	host := C.GoString(name)
	addrs, err := net.LookupHost(host)
	if err != nil {
		log.Fatalf("Unable to lookup host %s: %s", host, err)
	}

	addrIP := make([]C.struct_in_addr, len(addrs)+1)
	addrIP[len(addrIP)-1].s_addr = C.INADDR_NONE

	for index, addr := range addrs {
		C.inet_pton(C.AF_INET, C.CString(addr), unsafe.Pointer(&addrIP[index]))
	}
	_, err = syscall.Write((int)(fd), C.GoBytes(unsafe.Pointer(&addrIP[0]), (C.int)((len(addrIP))*C.sizeof_struct_in_addr)))
	if err != nil {
		log.Fatalf("Unable to send address: %s", err)
	}
	syscall.Close((int)(fd))
	return 0
}

//export cleanupGoResolver
func cleanupGoResolver() {
}

func init() {
	C.SetCustomResolver()
}
