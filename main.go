package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	i3ipcHeader = "i3-ipc"
)

func main() {
	conn, err := connectToIPC()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		return
	}

	// Subscribe to events
	if err := subscribeToEvents(conn, []string{"window"}); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		return
	}

	eventCh := make(chan *i3Message)
	msgCh := make(chan *i3Message)

	go handleMessageChannel(msgCh)
	go handleEventChannel(eventCh)

	handleApplicationRead(msgCh, eventCh, conn)
}

func handleApplicationRead(msgCh, eventCh chan *i3Message, conn net.Conn) {
	buff := make([]byte, 2024)

	for {
		n, err := conn.Read(buff[:cap(buff)])
		if err != nil {
			break
		}
		buff = buff[:n]

		// Header information
		if string(buff[0:len(i3ipcHeader)]) != i3ipcHeader {
			continue
		}

		// Message length
		messageLength := binary.LittleEndian.Uint32(buff[len(i3ipcHeader) : len(i3ipcHeader)+4])

		// Message type
		messageType := binary.LittleEndian.Uint32(buff[len(i3ipcHeader)+4 : len(i3ipcHeader)+8])
		messageTypeCategory := messageType & 0x7F

		messageData := buff[len(i3ipcHeader)+8:]

		i3Message := &i3Message{
			data:     messageData,
			category: messageTypeCategory,
			length:   messageLength,
		}

		// Message is not an event
		if messageType>>31 != 1 {
			msgCh <- i3Message
			continue
		}
		eventCh <- i3Message
	}
}

func handleMessageChannel(ch chan *i3Message) {
	for {
		select {
		case msg := <-ch:
			fmt.Println("msg", msg)
		}
	}
}

func handleEventChannel(ch chan *i3Message) {
	for {
		select {
		case msg := <-ch:
			fmt.Println("event", msg)
		}
	}
}
