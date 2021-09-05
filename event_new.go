package main

import (
	"fmt"
	"net"
)

type eventNewHandler struct {
	conn  net.Conn
	chans *eventHandlerChannels
}

func (e *eventNewHandler) handle(msg *i3Message, event *i3MessageEvent) error {
	tree, err := retrievei3Tree(e.chans.getTreeCh, e.conn)
	if err != nil {
		return err
	}

	treeData, err := tree.decodeTreeData()
	if err != nil {
		return err
	}

	fmt.Println(treeData.retrieveWorkspacesInformation())
	return nil
}
