package hub

import (
	"fmt"
	"net"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.uber.org/zap"
)

type EventHandler func(*Hub, *Client, Event)

type Hub struct {
	Clients    map[string]*Client
	unregister chan string
	broadcast  chan []byte
	OnEvent    EventHandler
	wg         sync.WaitGroup
	// id         uint64
}

func newID() (string, error) {
	return gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 6)
}

func NewEventID() string {
	return fmt.Sprintf("evt_%s", gonanoid.MustGenerate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 6))
}

func NewHub() (*Hub, error) {
	return &Hub{
		Clients:    make(map[string]*Client),
		unregister: make(chan string),
		broadcast:  make(chan []byte),
	}, nil
}

func (h *Hub) Broadcast(evt Event) {
	b, _ := json.Marshal(evt)
	h.broadcast <- b
}

func (h *Hub) SendTo(id string, evt Event) error {
	if client, ok := h.Clients[id]; ok {
		b, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		client.out <- b
	}
	return nil
}

func (h *Hub) RegisterClient(conn net.Conn) string {
	id, _ := newID()
	h.Clients[id] = &Client{
		conn:      conn,
		SessionID: id,
		out:       make(chan []byte, 10),
		hub:       h,
	}
	defer h.Clients[id].Start()
	return id
}

func (h *Hub) Listen() {
	defer h.wg.Done()
	for {
		select {
		case id := <-h.unregister:
			if _, ok := h.Clients[id]; ok {
				// close(client.out)
				delete(h.Clients, id)
				Logger.Debug("client remove",
					zap.String("id", id),
				)
			}
		case message := <-h.broadcast:
			for _, client := range h.Clients {
				client.out <- message
			}
		}
	}
}

func (h *Hub) Run() {
	Logger.Info("chat-server start",
		zap.Time("start at", time.Now()),
	)
	h.wg.Add(1)
	go h.Listen()
	h.wg.Wait()
}
