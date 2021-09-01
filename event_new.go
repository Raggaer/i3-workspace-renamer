package main

import (
	"net"
)

type eventNewHandler struct {
	conn  net.Conn
	chans *eventHandlerChannels
}

func (e *eventNewHandler) handle(msg *i3Message, event *i3MessageEvent) error {
	return nil
}
