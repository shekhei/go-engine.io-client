package main

import (
  "github.com/gorilla/websocket"
  "log"
  "encoding/json"
)

type Event struct {
  Type string
  Data []byte
}

type Client struct {
  Conn *websocket.Conn
  sender chan *Packet
  receiver chan *Packet
  event chan *Event
  done chan bool
  opened bool
}

func NewClient(conn *websocket.Conn) (*Client) {
  return &Client{
    Conn: conn,
    opened: false,
    sender: make(chan *Packet),
    receiver: make(chan *Packet),
    event: make(chan *Event),
    done: make(chan bool),
  }
}

func Dial(urlStr string) (client *Client, err error) {
  urlStr += "/engine.io/?transport=websocket"
  log.Println(urlStr)
  conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
  if nil != err {
    return
  }
  client = NewClient(conn)
  go func(){
    // listens for open event, and then make it open
    for {
      _, res, err := client.Conn.ReadMessage()
      if nil != err {
        panic(err)
      }
      log.Printf("%s", string(res))
      packet := BytesToPacket(res)
      client.receiver <- packet
    }
  }()
  go func() {
    for {
      packet := <- client.receiver
      event := &Event{Data: packet.Data}
      switch packet.Type {
        case Open:
          client.opened = true
          event.Type = "open"
        case Close:
          client.opened = false
          event.Type = "close"
        case Message:
          event.Type = "message"
        case Ping:
          event.Type = "ping"
          log.Println("Got a ping event, replying now")
          client.sender <- &Packet{Type:Pong, Data:packet.Data}
        case Pong:
          event.Type = "pong"
        default:
          event = nil
      }
      if nil != event {
        client.event <- event
      }
    }
  }()
  go func() {
    for {
      select {
        case p := <-client.sender:
          err := client.Conn.WriteMessage(websocket.TextMessage, PacketToBytes(p))
          if nil != err {
            panic(err)
          }
        default:
      }
    }
  }()
  return
}

func (c *Client) SendMessage(obj interface{}) {
  result, err := json.Marshal(obj)
  if nil != err {
    panic(err)
  }
  c.sender <- NewPacket(Message, result)
}

func (c *Client) SendPacket(packet *Packet) {
  c.sender <- packet
}

func (c *Client) Close() error {
  return c.Conn.Close()
}
