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

const (
	GG_CHECK_READ  = C.GG_CHECK_READ
	GG_CHECK_WRITE = C.GG_CHECK_WRITE
)

type GGSession struct {
	Uin      uint32
	Password string
	session  *C.struct_gg_session

	events chan *GGEvent
}

func NewGGSession() *GGSession {
	s := new(GGSession)
	s.session = nil
	s.events = make(chan *GGEvent, 100)
	return s
}

func Version() string {

	return C.GoString(C.gg_libgadu_version())
}

func (session GGSession) Close() {
	C.gg_free_session(session.session)
}

func (session GGSession) watchFd() *GGEvent {
	return (*GGEvent)(C.gg_watch_fd(session.session))
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
	go session.poller()
	for {
		e := <-session.events
		defer e.Close()

		if e.Type() == GG_EVENT_CONN_FAILED {
			return AccessDeniedError
		}
		if e.Type() == GG_EVENT_CONN_SUCCESS {
			return nil
		}
	}
	return Fault
}
