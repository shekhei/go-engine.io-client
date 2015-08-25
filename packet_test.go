package main

import (
  "testing"
  "encoding/json"
)

type TestStruct struct {
  Str string
  Num int64
  Arr []string
}

func TestPacketEncoding(test* testing.T) {
  packet := NewEmptyPacket(Open)
  if Open != packet.Type {
    test.Errorf("Expecting packet to be equal to %s and not %s", string(packet.Type), string(Open))
  }
  if 0 != len(packet.Data) {
    test.Errorf("Expecting packet data to be empty")
  }
  newStruct := &TestStruct{"abc", 12, []string{"def", "cde"}}
  result, err := json.Marshal(newStruct)
  if nil != err {
    test.Errorf("Error is thrown %v", err)
  }
  packet = NewPacket(Message, result)

  if Message != packet.Type {
    test.Errorf("Expecting packet to be equal to %s and not %s", string(packet.Type), string(Message))
  }
  if len(result) != len(packet.Data) {
    test.Errorf("Expecting packet data to be the same as json marshalled")
  }

  encoded := PacketToBytes(packet)
  decoded := BytesToPacket(encoded)
  test.Logf("Lets see how it looks like %s", string(encoded))
  if Message != packet.Type {
    test.Errorf("Expecting packet to be equal to %s and not %s", string(packet.Type), string(Message))
  }
  if len(result) != len(decoded.Data) {
    test.Errorf("Expecting packet data to be the same as json marshalled")
  }
  struct2 := new(TestStruct)
  json.Unmarshal(decoded.Data, struct2)
  if newStruct.Str != struct2.Str {
    test.Errorf("Expecting struct2.Str to be equal to %s and not %s", string(struct2.Str), string(newStruct.Str))
  }
  if newStruct.Num != struct2.Num {
    test.Errorf("Expecting struct2.Num to be equal to %s and not %s", string(struct2.Num), string(newStruct.Num))
  }
  if len(newStruct.Arr) != len(struct2.Arr) {
    test.Errorf("Expecting the len(struct2.Arr) to equal len(newStruct.Arr)")
  }
  if newStruct.Arr[0] != struct2.Arr[0] ||
      newStruct.Arr[1] != struct2.Arr[1] {
    test.Errorf("Expecting the struct2.Arr to equal newStruct.Arr")
  }

}
