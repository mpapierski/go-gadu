// Package gadu implements a wrapper over asynchronous interface in libgadu
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
	// GGFeatureImageDescr means that client supports graphical statuses
	GGFeatureImageDescr = C.GG_FEATURE_IMAGE_DESCR
)

const (
	ggCheckRead  = C.GG_CHECK_READ
	ggCheckWrite = C.GG_CHECK_WRITE
)

// GGSession is a struct representing session with GG network
type GGSession struct {
	Uin      uint32
	Password string
	session  *C.struct_gg_session

	events chan *GGEvent
}

// NewGGSession returns new instance of a session
func NewGGSession() *GGSession {
	s := new(GGSession)
	s.session = nil
	s.events = make(chan *GGEvent, 100)
	return s
}

// Version returns a current version of the library
func Version() string {
	return C.GoString(C.gg_libgadu_version())
}

// Close frees resources allocated by the library
func (session GGSession) Close() {
	C.gg_free_session(session.session)
	session.session = nil
}

func (session GGSession) watchFd() *GGEvent {
	return (*GGEvent)(C.gg_watch_fd(session.session))
}

// Login starts connection with the server
func (session GGSession) Login() error {
	params := C.struct_gg_login_params{
		uin:               (C.uin_t)(session.Uin),
		password:          C.CString(session.Password),
		async:             (C.int)(1),
		protocol_features: GGFeatureImageDescr,
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

		if e.Type() == GGEventConnectionFailed {
			return AccessDeniedError
		}
		if e.Type() == GGEventConnectionSuccess {
			return nil
		}
	}
	return Fault
}

// Logout ends connection with the server
func (session GGSession) Logout() {
	C.gg_logoff(session.session)
}
