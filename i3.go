package main

import (
	"encoding/json"
	"net"
)

type i3Message struct {
	length   uint32
	category uint32
	isEvent  bool
	data     []byte
}

func retrievei3Workspaces(conn net.Conn) error {
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
