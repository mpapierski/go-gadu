package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>
*/
import "C"

// GGEvent is a sugar for C struct of the same meaning
type GGEvent C.struct_gg_event

const (
	// GGEventConnectionSuccess when established connection
	GGEventConnectionSuccess = C.GG_EVENT_CONN_SUCCESS
	// GGEventConnectionFailed when failed to establish connection
	GGEventConnectionFailed = C.GG_EVENT_CONN_FAILED
)

// Close frees memory allocated with event
func (event *GGEvent) Close() {
	C.gg_free_event((*C.struct_gg_event)(event))
}

// Type gets information about internal type of the event
func (event *GGEvent) Type() int {
	return (int)(event._type)
}
