package main

import (
	"net"
)

const (
	EVENT_NEW   = "new"
	EVENT_CLOSE = "close"
)

type eventHandlerChannels struct {
	getWorkspacesCh chan *i3Message
	getTreeCh       chan *i3Message
}

type eventHandler interface {
	handle(*i3Message, *i3MessageEvent) error
}

func handleMessageChannel(ch, getWorkspacesCh, getTreeCh chan *i3Message) {
	for {
		select {
		case msg := <-ch:
			switch msg.category {
			case I3_CATEGORY_GET_TREE:
				getTreeCh <- msg
			case I3_CATEGORY_GET_WORKSPACES:
				getWorkspacesCh <- msg
			default:
				continue
			}
		}
	}
}

func handleEventChannel(conn net.Conn, ch, getWorkspacesCh, getTreeCh chan *i3Message) {
	channelHolder := &eventHandlerChannels{
		getWorkspacesCh: getWorkspacesCh,
		getTreeCh:       getTreeCh,
	}

	for {
		select {
		case msg := <-ch:
			// Decode event data
			data, err := msg.decodeEventData()
			if err != nil {
				continue
			}

			// Handle event
			var handler eventHandler
			switch data.Change {
			case EVENT_NEW, EVENT_CLOSE:
				handler = &eventNewHandler{conn, channelHolder}
				go handler.handle(msg, data)
			}
		}
	}
}
