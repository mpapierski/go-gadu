package gadu

/*
#cgo pkg-config: libgadu
#include <libgadu.h>
*/
import "C"

type GGEvent C.struct_gg_event

const (
	GG_EVENT_CONN_SUCCESS = C.GG_EVENT_CONN_SUCCESS
	GG_EVENT_CONN_FAILED  = C.GG_EVENT_CONN_FAILED
)

func (event *GGEvent) Close() {
	C.gg_free_event((*C.struct_gg_event)(event))
}

func (event *GGEvent) Type() int {
	return (int)(event._type)
}
