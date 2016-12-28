package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>
#include <errno.h>

static struct gg_session* gg_login2(struct gg_login_params* params, int* error) {
	struct gg_session* session = gg_login(params);
	if (!session) {
		*error = errno;
	}
	return session;
}
*/
import "C"
import "unsafe"

const (
	GG_FEATURE_IMAGE_DESCR = C.GG_FEATURE_IMAGE_DESCR
)

type GGSession struct {
	Uin      uint32
	Password string
	session  *C.struct_gg_session
}

func NewGGSession() *GGSession {
	s := new(GGSession)
	s.session = nil
	return s
}

func Version() string {

	return C.GoString(C.gg_libgadu_version())
}

func (session GGSession) Login() error {
	params := C.struct_gg_login_params{
		uin:               (C.uin_t)(session.Uin),
		password:          C.CString(session.Password),
		async:             (C.int)(1),
		protocol_features: GG_FEATURE_IMAGE_DESCR,
	}

	// We use wrapper because a wrapper function because cgo doesnt like errno.
	// See https://github.com/golang/go/issues/1360
	var err int
	session.session = C.gg_login2((*C.struct_gg_login_params)(unsafe.Pointer(&params)), (*C.int)(unsafe.Pointer(&err)))
	if session.session == nil {
		return NewGGError(err)
	}
	return nil
}
