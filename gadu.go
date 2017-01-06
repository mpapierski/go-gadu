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

int gg_ping2(struct gg_session* session) {
	if (gg_ping(session) == -1) {
		return -errno;
	}
	return 0;
}

int gg_send_message2(struct gg_session *sess, int msgclass, uin_t recipient, const char *message) {
	if (gg_send_message(sess, msgclass, recipient, (const unsigned char*)message) == -1) {
		return -errno;
	}
	return 0;
}

int gg_notify2(struct gg_session *sess, uin_t *userlist, int count) {
	if (gg_notify(sess, userlist, count) == -1) {
		return -errno;
	}
	return 0;
}
*/
import "C"
import (
	"log"
	"time"
	"unsafe"
)

const (
	// GGFeatureImageDescr means that client supports graphical statuses
	GGFeatureImageDescr = C.GG_FEATURE_IMAGE_DESCR
)

const (
	ggCheckRead  = C.GG_CHECK_READ
	ggCheckWrite = C.GG_CHECK_WRITE
)

const (
	ggClassMsg = C.GG_CLASS_MSG
)

// Uin is your GG number
type Uin C.uin_t

// GGSession is a struct representing session with GG network
type GGSession struct {
	Uin        Uin
	Password   string
	session    *C.struct_gg_session
	Events     chan *GGEvent
	pingTicker *time.Ticker
}

// NewGGSession returns new instance of a session
func NewGGSession() *GGSession {
	s := new(GGSession)
	s.session = nil
	s.Events = make(chan *GGEvent, 100)
	s.pingTicker = nil
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
	if session.pingTicker != nil {
		session.pingTicker.Stop()
	}
}

func (session GGSession) watchFd() *GGEvent {
	return &GGEvent{cptr: C.gg_watch_fd(session.session)}
}

// Login starts connection with the server
func (session *GGSession) Login() error {
	params := C.struct_gg_login_params{
		uin:               (C.uin_t)(session.Uin),
		password:          C.CString(session.Password),
		async:             (C.int)(1),
		protocol_features: GGFeatureImageDescr,
	}

	// We use wrapper because a wrapper function because cgo doesnt like errno.
	// See https://github.com/golang/go/issues/1360
	var err int

	s := C.gg_login2((*C.struct_gg_login_params)(unsafe.Pointer(&params)), (*C.int)(unsafe.Pointer(&err)))
	if s == nil {
		return NewGGError(err)
	}
	session.session = s
	go session.poller()
	for {
		e := <-session.Events
		defer e.Close()

		if e.Type() == GGEventConnectionFailed {
			return AccessDeniedError
		}
		if e.Type() == GGEventConnectionSuccess {
			// Start ping timer
			session.pingTicker = time.NewTicker(60 * time.Second)
			go session.ticker()
			return nil
		}
	}
	return Fault
}

func (session GGSession) ping() error {
	if result := C.gg_ping2(session.session); result != 0 {
		return NewGGError(-(int)(result))
	}
	return nil
}

func (session GGSession) ticker() {
	for range session.pingTicker.C {
		if err := session.ping(); err != nil {
			log.Fatalf("Unable to send ping: %s", err)
			break
		}
	}
}

// Logout ends connection with the server
func (session GGSession) Logout() {
	C.gg_logoff(session.session)
}

// Notify sends contact list to server
func (session GGSession) Notify(userList []Uin) error {
	if errno := C.gg_notify2(session.session, (*C.uin_t)(unsafe.Pointer(&userList[0])), (C.int)(len(userList))); errno != 0 {
		return NewGGError(-(int)(errno))
	}
	return nil
}

func (session GGSession) GetCPtr() *C.struct_gg_session {
	return session.session
}

// SendMessage sends a message to a recipient
func (session GGSession) SendMessage(uin Uin, text string) error {
	if errno := C.gg_send_message2(session.session, (C.int)(ggClassMsg), (C.uin_t)(uin), C.CString(text)); errno != 0 {
		return NewGGError(-(int)(errno))
	}
	return nil
}
