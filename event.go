package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>

struct gg_event_msg gg_event_get_msg(struct gg_event* event) {
	return event->event.msg;
}

*/
import "C"
import "unsafe"

// GGEvent is a sugar for C struct of the same meaning
type GGEvent struct {
	cptr *C.struct_gg_event
}

// GGEventMsg is a sugar for message event
type GGEventMsg C.struct_gg_event_msg

const (
	// GGEventConnectionSuccess when established connection
	GGEventConnectionSuccess = C.GG_EVENT_CONN_SUCCESS
	// GGEventConnectionFailed when failed to establish connection
	GGEventConnectionFailed = C.GG_EVENT_CONN_FAILED
	// GGEventMsg when received message
	GGEventMessage = C.GG_EVENT_MSG
)

// Close frees memory allocated with event
func (event *GGEvent) Close() {
	C.gg_free_event((*C.struct_gg_event)(event.cptr))
}

// Type gets information about internal type of the event
func (event *GGEvent) Type() int {
	return (int)(event.cptr._type)
}

// Message gets message event details
func (event *GGEvent) Message() GGEventMsg {
	return (GGEventMsg)(C.gg_event_get_msg(event.cptr))
}

func (event *GGEventMsg) Sender() Uin {
	return (Uin)(event.sender)
}

func (event *GGEventMsg) Message() string {
	return C.GoString((*C.char)(unsafe.Pointer(event.message)))
}
