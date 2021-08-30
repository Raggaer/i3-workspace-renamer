package main

import (
	"encoding/binary"
)

type outputMessage struct {
	pos  int
	data []byte
}

func newOutputMessage() *outputMessage {
	return &outputMessage{
		pos:  0,
		data: make([]byte, 500),
	}
}

func (o *outputMessage) writeUint32(v uint32) {
	binary.LittleEndian.PutUint32(o.data[o.pos:], v)
	o.pos += 4
}

func (o *outputMessage) writeString(v string) {
	copy(o.data[o.pos:], v)
	o.pos += len(v)
}

func (o *outputMessage) writeBytes(v []byte) {
	copy(o.data[o.pos:], v)
	o.pos += len(v)
}

func (o *outputMessage) getOutputData() []byte {
	return o.data[0:o.pos]
}
