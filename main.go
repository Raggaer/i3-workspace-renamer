package main

import (
	"encoding/binary"
	"flag"
	"net"
	"time"
)

const (
	i3ipcHeader = "i3-ipc"
	bufferSize  = 61440
)

func main() {
	// Retrieve config file path
	var cfgFilePath = flag.String("config", "i3wr_config.json", "Filepath for the application config file")
	flag.Parse()

	eventCh := make(chan *i3Message)
	msgCh := make(chan *i3Message)
	getWorkspacesCh := make(chan *i3Message)
	getTreeCh := make(chan *i3Message)

	go handleMessageChannel(msgCh, getWorkspacesCh, getTreeCh)

	retryStartApplication(*cfgFilePath, getWorkspacesCh, getTreeCh, msgCh, eventCh)
}

func retryStartApplication(configFilePath string, getWorkspacesCh, getTreeCh, msgCh, eventCh chan *i3Message) {
	for {
		conn, err := connectToIPC()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		// Subscribe to events
		if err := subscribeToEvents(conn, []string{"window", "workspace", "output", "binding", "shutdown"}); err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		// Load config file
		cfg, err := loadConfigurationFile(configFilePath)
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		go handleEventChannel(cfg, conn, eventCh, getWorkspacesCh, getTreeCh)

		// Handle application reading
		handleApplicationRead(msgCh, eventCh, conn)

		// If application stopped reading reset
		eventCh <- nil
		time.Sleep(time.Second * 5)
	}
}

func handleApplicationRead(msgCh, eventCh chan *i3Message, conn net.Conn) error {
	buff := make([]byte, bufferSize)

	for {
		n, err := conn.Read(buff[:cap(buff)])
		if err != nil {
			return err
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

		messageData := make([]byte, len(buff[len(i3ipcHeader)+8:]))
		copy(messageData, buff[len(i3ipcHeader)+8:])

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
