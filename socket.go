package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func getSocketPath() (string, error) {
	out, err := exec.Command("i3", "--get-socketpath").Output()
	return string(out), err
}

func connectToIPC() (net.Conn, error) {
	path, err := getSocketPath()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve i3 socketpath: %v", err)
	}

	conn, err := net.Dial("unix", strings.TrimSpace(path))
	return conn, err
}
