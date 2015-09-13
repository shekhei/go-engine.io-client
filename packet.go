package engineio

import (
  // "encoding/json"
  "log"
)

type PacketType byte

type Packet struct {
  Type PacketType
  Data []byte
}

const (
  Open PacketType = '0'
  Close PacketType = '1'
  Ping PacketType = '2'
  Pong PacketType = '3'
  Message PacketType = '4'
  Upgrade PacketType = '5'
  Noop PacketType = '6'
)

func NewEmptyPacket(pType PacketType) (*Packet) {
  return &Packet{pType, make([]byte, 0)}
}

func NewPacket(pType PacketType, content []byte) (*Packet) {
  return &Packet{pType, content}
}

func NewClosePacket() (*Packet) {
  return &Packet{Close, nil}
}

func BytesToPacket(from []byte) (to *Packet) {
  return &Packet{Type: PacketType(from[0]), Data: from[1:]}
}

func PacketToBytes(from *Packet) (to []byte) {
  to = make([]byte, len(from.Data)+1)
  copy(to[1:], from.Data)
  to[0] = byte(from.Type)
  log.Println(string(to))
  return
}
