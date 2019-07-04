package main

import (
	"encoding/binary"
	"strconv"
)

// TypeEnum enum
type TypeEnum int32

const (
	// Ack type
	Ack TypeEnum = 0
	// Data type
	Data TypeEnum = 1
	// Eot type
	Eot TypeEnum = 2
)

// Packet struct
type Packet struct {
	Type   TypeEnum
	Seqnum int32
	Data   []byte
}

// GetBytes func
func (p *Packet) GetBytes() []byte {
	bs := make([]byte, 12, 512)
	binary.BigEndian.PutUint32(bs, uint32(p.Type))
	binary.BigEndian.PutUint32(bs[4:], uint32(p.Seqnum))
	binary.BigEndian.PutUint32(bs[8:], uint32(len(p.Data)))
	bs = append(bs, p.Data...)
	return bs
}

// NewPacket func
func NewPacket(bs []byte) (*Packet, error) {
	p := Packet{
		Type:   TypeEnum(binary.BigEndian.Uint32(bs[0:4])),
		Seqnum: int32(binary.BigEndian.Uint32(bs[4:8])),
	}
	plen := binary.BigEndian.Uint32(bs[8:12])
	p.Data = bs[12 : 12+plen]
	return &p, nil
}

func (p *Packet) String() string {
	return "Packet[type=" + strconv.Itoa(int(p.Type)) + "]" +
		"[seq=" + strconv.Itoa(int(p.Seqnum)) + "]" +
		"[len=" + strconv.Itoa(len(p.Data)) + "]"
}
