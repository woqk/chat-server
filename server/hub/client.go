package hub

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Client struct {
	hub       *Hub
	conn      net.Conn
	SessionID string
	UserID    string
	IsAuth    bool
	// Channels []string
	out chan []byte
}

func (c *Client) reader() error {
	defer func() {
		// if r := recover(); r != nil {
		// 	fmt.Println("Recovered in f", r)
		// }
		// c.out = nil
		c.Disconnect()
	}()
	for {
		header, err := ws.ReadHeader(c.conn)
		if err != nil {
			// handle error
			return err
		}

		payload := make([]byte, header.Length)
		_, err = io.ReadFull(c.conn, payload)
		if err != nil {
			// handle error
			return err
		}
		if header.Masked {
			ws.Cipher(payload, header.Mask, 0)
		}

		// // Reset the Masked flag, server frames must not be masked as
		// // RFC6455 says.
		header.Masked = false

		if header.OpCode == ws.OpClose {
			Logger.Debug("connection closed",
				zap.String("session_id", c.SessionID),
				zap.String("user_id", c.UserID),
			)
			break
		}

		evt := Event{}
		if err := json.Unmarshal(payload, &evt); err != nil {
			log.Fatal(err)
		}
		evt.SentAt = time.Now()

		evt.Meta = Meta{
			SessionID: c.SessionID,
			UserID:    c.UserID,
			// TODO: Add User
		}
		c.hub.OnEvent(c.hub, c, evt)
		// c.hub.Events <- evt
	}
	return nil
}

func (c *Client) writer() error {
	// fmt.Printf("Write work for: %s", c.SessionID)
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered in f", r)
	// 	}
	// }()
	// defer fmt.Printf("Stop writer(): %s\n", c.SessionID)
	for {
		select {
		case rep, ok := <-c.out:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				ws.WriteFrame(c.conn, ws.NewCloseFrame([]byte("")))
				return nil
			}

			frame := ws.NewTextFrame(rep)
			if err := ws.WriteFrame(c.conn, frame); err != nil {
				// handle err
				return err
			}
		}
	}
}

func (c *Client) Dispatch(evt Event) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	c.out <- b
	return nil
}

func (c *Client) Start() {
	go c.reader()
	go c.writer()
}

func (c *Client) Disconnect() {
	close(c.out)
	c.conn.Close()
	c.hub.unregister <- c.SessionID
}
