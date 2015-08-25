package main

import (
  "testing"
  "time"
  "encoding/json"
  "strconv"
)

var opened = false

func HandleNotOpen(test* testing.T, ev *Event, client *Client) {
  if !opened && ev.Type != "open" {
    test.Errorf("The first event should be open!")
  }
  opened = true
  client.SendPacket(&Packet{Type: Ping, Data: []byte("test")})
}

func HandleOpened(test* testing.T, ev *Event, client *Client) {
  if ev.Type != "pong" {
    test.Errorf("The next event should be pong!")
  }
}

func TestClient(test* testing.T) {
  client, err := Dial("ws://localhost:55555")
  // defer client.Close()
  if nil != err {
    test.Errorf("it failed toc onnect?! %s", err.Error())
  }
  Loop:
    for {
      select {
        case ev := <- client.event:
          if !opened {
            HandleNotOpen(test, ev, client)
          } else {
            HandleOpened(test, ev, client)
            break Loop
          }
        default:
      }
    }
}

type SendingObj struct {
  Type string
  Msg string
}

func TestClientMessaging(test* testing.T) {
  opened = false
  client, err := Dial("ws://localhost:55555")
  // defer client.Close()
  if nil != err {
    test.Errorf("it failed toc onnect?! %s", err.Error())
  }
  timestamp := strconv.FormatInt(time.Now().Unix(), 10)
  Loop:
    for {
      select {
        case ev := <- client.event:
          if "open" == ev.Type {
            client.SendMessage(&SendingObj{Type: "echo", Msg: timestamp})
          } else if "message" == ev.Type {
            tempStruct := new(SendingObj)
            json.Unmarshal(ev.Data, tempStruct)
            if "echo" != tempStruct.Type {
              test.Errorf("This should return as an echo type")
            }
            if timestamp+"abc" != tempStruct.Msg {
              test.Errorf("Expecting timestampabc and not '%s'", tempStruct.Msg)
            }
            break Loop
          }
        default:
      }
    }
}
