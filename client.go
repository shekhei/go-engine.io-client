package engineio

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type string
	Data []byte
}

var closeErrorType = reflect.TypeOf((*websocket.CloseError)(nil))

type Client struct {
	Conn     *websocket.Conn
	sender   chan *Packet
	receiver chan *Packet
	Event    chan *Event
	done     chan bool
	opened   bool
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn:     conn,
		opened:   false,
		sender:   make(chan *Packet, 5),
		receiver: make(chan *Packet, 5),
		Event:    make(chan *Event, 5),
		done:     make(chan bool),
	}
}

func Dial(urlStr string) (client *Client, err error) {
	urlStr += "/engine.io/?transport=websocket"
	log.Println(urlStr)
	client = NewClient(nil)
	go func() {
		// listens for open event, and then make it open
		for {
			if nil == client.Conn {
				continue
			}
			_, res, err := client.Conn.ReadMessage()
			log.Println("Reading next Message %s, %s", res, err)
			if nil != err {
				t := reflect.TypeOf(err)
				log.Println("error, %T(%q)", t, err)
				if closeErrorType == t {
					client.Conn = nil
					// handle reconnections
					client.receiver <- NewClosePacket()
				} else {
					panic(err)
				}

			}
			log.Printf("%s", string(res))
			packet := BytesToPacket(res)
			log.Println("Decoding packet %q", packet.Type)
			client.receiver <- packet
		}
	}()
	go func() {
		for {
			packet := <-client.receiver
			log.Println("Got a packet!!! yeah!!")
			event := &Event{Data: packet.Data}
			switch packet.Type {
			case Open:
				client.opened = true
				event.Type = "open"
			case Close:
				client.opened = false
				event.Type = "close"
			case Message:
				log.Println("Got a message event!!! yeah!!")
				event.Type = "message"
			case Ping:
				event.Type = "ping"
				log.Println("Got a ping event, replying now")
				client.sender <- &Packet{Type: Pong, Data: packet.Data}
			case Pong:
				event.Type = "pong"
			default:
				event = nil
			}
			if nil != event {
				log.Println("Sending event out!")
				client.Event <- event
			}
		}
	}()
	go func() {
		// reconnection ticker
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case _ = <-ticker.C:
				// do reconnection here if client is not connected
				// log.Println("reconnecting ticker")
				if nil == client.Conn {
					log.Println("reconnecting")
					conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
					if nil != err {
						log.Println("reconnecting fucked up %q", err)
						break
					}
					log.Println("reconnecting succeeded")
					client.Conn = conn
				}
				time.Sleep(time.Second)
			default:
			}
		}
	}()
	go func() {
		for {
			select {
			case p := <-client.sender:
				if client.Conn == nil {
					break
				}
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
	log.Printf(string(result))
	c.sender <- NewPacket(Message, result)
}

func (c *Client) SendPacket(packet *Packet) {
	c.sender <- packet
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
