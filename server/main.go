package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/gobwas/ws"
	"github.com/imdario/mergo"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/thoas/go-funk"
	serv "github.com/woqk/chat-server/server/hub"
	"go.uber.org/zap"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type TempDB struct {
	Channels     []string
	UserSessions map[string][]string
	UserChannels map[string][]string
	Users        map[string]User
}

var (
	tempDB *TempDB
)

func init() {
	tempDB = &TempDB{
		Channels: []string{
			"public",
			"design",
			"go",
			"svelte",
		},
		Users:        make(map[string]User),
		UserSessions: make(map[string][]string),
		UserChannels: make(map[string][]string),
	}
}

func SendToChannel(hub *serv.Hub, channel string, evt serv.Event) error {
	for _, client := range hub.Clients {
		if funk.ContainsString(tempDB.UserChannels[client.UserID], channel) {
			client.Dispatch(evt)
		}
	}
	return nil
}

func OnEvent(hub *serv.Hub, client *serv.Client, evt serv.Event) {
	if evt.Type == serv.TypeLogin {
		client.IsAuth = true
		client.UserID = fmt.Sprintf("usr_%s", gonanoid.Must(5))
		tempDB.Users[client.UserID] = User{
			ID:       client.UserID,
			Username: evt.Body["username"].(string),
		}
		client.Dispatch(serv.Event{
			ID:     serv.NewEventID(),
			Type:   "login",
			SentAt: evt.SentAt,
			Body: map[string]interface{}{
				"ok":        true,
				"username":  evt.Body["username"],
				"userId":    client.UserID,
				"sessionId": client.SessionID,
			},
		})
		tempDB.UserSessions[client.UserID] = []string{client.SessionID}
		return
	}

	if !client.IsAuth {
		return
	}

	switch evt.Type {
	case serv.TypeChannelJoin:
		channel, ok := evt.Body["channel"].(string)
		if !ok {
			// TODO: Send error
			return
		}

		joinedChannel, ok := tempDB.UserChannels[evt.Meta.UserID]
		if !ok {
			tempDB.UserChannels[evt.Meta.UserID] = []string{}
		}

		if funk.ContainsString(joinedChannel, channel) {
			client.Dispatch(serv.CreateErrorEvent(-1, fmt.Sprintf("Already joined channel \"%s\"", channel)))
			serv.Logger.Debug("user already joinned channel",
				zap.String("user_id", evt.Meta.UserID),
				zap.String("session_id", evt.Meta.SessionID),
				zap.String("channel", channel),
			)
			return
		}
		tempDB.UserChannels[evt.Meta.UserID] = append(tempDB.UserChannels[evt.Meta.UserID], channel)
		serv.Logger.Debug("user join channel",
			zap.String("user_id", evt.Meta.UserID),
			zap.String("session_id", evt.Meta.SessionID),
			zap.String("channel", channel),
		)
		SendToChannel(hub, channel, serv.Event{
			ID:     serv.NewEventID(),
			Type:   serv.TypeChannelJoined,
			SentAt: evt.SentAt,
			Body: map[string]interface{}{
				"channel": channel,
				"sender":  tempDB.Users[client.UserID],
				// "user": map[string]interface{}{
				// 	"id": evt.Meta.UserID,
				// },
			},
		})

	case serv.TypeChannelList:
		client.Dispatch(serv.Event{
			ID:     serv.NewEventID(),
			Type:   serv.TypeChannelResult,
			SentAt: evt.SentAt,
			Body: map[string]interface{}{
				"channels": tempDB.Channels,
			},
		})

	case serv.TypeFriendList:
		client.Dispatch(serv.Event{
			ID:     serv.NewEventID(),
			Type:   serv.TypeFriendResult,
			SentAt: evt.SentAt,
			Body: map[string]interface{}{
				"friends": tempDB.Users,
			},
		})

	case serv.TypePostMessage:
		body := map[string]interface{}{
			// "channel":  evt.Body["channel"],
			// "user":     evt.Body["user"],
			// "replyTo": evt.Body["replyTo"],
			// "text":     evt.Body["text"],
		}

		if err := mergo.Merge(&body, evt.Body, mergo.WithOverride); err != nil {
			log.Fatal(err)
		}
		body["sender"] = tempDB.Users[client.UserID]

		sEvt := serv.Event{
			ID:     serv.NewEventID(),
			Type:   serv.TypeMessage,
			SentAt: evt.SentAt,
			Body:   body,
		}

		if user, ok := evt.Body["user"].(string); ok {
			// TODO: Check is user friend
			hub.SendTo(user, sEvt)
		} else if channel, ok := evt.Body["channel"].(string); ok {
			// TODO: Check is user join channel
			SendToChannel(hub, channel, sEvt)
		} else {
			hub.Broadcast(sEvt)
		}

	}
}

func main() {
	defer serv.Logger.Sync()
	hub, err := serv.NewHub()

	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	// Prepare handshake header writer from http.Header mapping.
	header := ws.HandshakeHeaderHTTP(http.Header{
		"X-Go-Version": []string{runtime.Version()},
	})
	hub.OnEvent = OnEvent
	go hub.Run()

	u := ws.Upgrader{
		// OnHost: func(host []byte) error {
		// 	if string(host) == "github.com" {
		// 		return nil
		// 	}
		// 	return ws.RejectConnectionError(
		// 		ws.RejectionStatus(403),
		// 		ws.RejectionHeader(ws.HandshakeHeaderString(
		// 			"X-Want-Host: github.com\r\n",
		// 		)),
		// 	)
		// },
		// OnHeader: func(key, value []byte) error {
		// 	if string(key) != "Cookie" {
		// 		return nil
		// 	}
		// 	ok := httphead.ScanCookie(value, func(key, value []byte) bool {
		// 		// Check session here or do some other stuff with cookies.
		// 		// Maybe copy some values for future use.
		// 		return true
		// 	})
		// 	if ok {
		// 		return nil
		// 	}
		// 	return ws.RejectConnectionError(
		// 		ws.RejectionReason("bad cookie"),
		// 		ws.RejectionStatus(400),
		// 	)
		// },
		OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {
			return header, nil
		},
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Fatal(err)
		}

		_, err = u.Upgrade(conn)
		if err != nil {
			// handle error
			log.Printf("upgrade error: %s", err)

		}
		id := hub.RegisterClient(conn)
		serv.Logger.Debug("New client",
			zap.String("id", id),
		)
	}

}
