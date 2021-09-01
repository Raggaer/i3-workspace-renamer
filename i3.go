package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	I3_CATEGORY_GET_WORKSPACES = 1
)

type i3MessageEvent struct {
	Change string `json:"change"`
}

type i3Message struct {
	length   uint32
	category uint32
	isEvent  bool
	data     []byte
}

func (i *i3Message) decodeEventData() (*i3MessageEvent, error) {
	var data i3MessageEvent
	err := json.Unmarshal(i.data, &data)
	return &data, err
}

func retrievei3Workspaces(ch chan *i3Message, conn net.Conn) (*i3Message, error) {
	if err := sendi3Workspaces(conn); err != nil {
		return nil, fmt.Errorf("Unable to retrieve i3 workspaces: %v", err)
	}
	msg := <-ch
	return msg, nil
}

func sendi3Workspaces(conn net.Conn) error {
	output := newOutputMessage()
	output.writeString("i3-ipc")
	output.writeUint32(0)
	output.writeUint32(1)

	data := output.getOutputData()
	_, err := conn.Write(data)
	return err
}

func subscribeToEvents(conn net.Conn, events []string) error {
	output := newOutputMessage()
	output.writeString("i3-ipc")

	subData, err := json.Marshal(events)
	if err != nil {
		return err
	}

	output.writeUint32(uint32(len(subData)))
	output.writeUint32(2)
	output.writeBytes(subData)

	data := output.getOutputData()
	_, err = conn.Write(data)
	return err
}
