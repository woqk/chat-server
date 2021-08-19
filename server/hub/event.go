package hub

import "time"

type Event struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Body   map[string]interface{} `json:"body,omitempty"`
	SentAt time.Time              `json:"sentAt"`
	Meta   Meta                   `json:"-"` // Server only data
}

func CreateErrorEvent(code int, message string) Event {
	return Event{
		Type: TypeError,
		Body: map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

type Meta struct {
	SessionID string
	UserID    string
}

// type EventLogin struct {
// 	Type string `json:"type"`
// 	Body struct {
// 		Username string `json:"username,omitempty"`
// 	}
// }

// type EventLoginResult struct {
// 	Type string `json:"type"`
// 	Body struct {
// 		OK bool `json:"ok"`
// 	}
// }

// type EventPostMessage struct {
// 	Type string `json:"type"`
// 	Body struct {
// 		ReplyTo *string `json:"replyTo,omitempty"`
// 		Channel *string `json:"channel,omitempty"`
// 		User    *string `json:"user,omitempty"`
// 		Text    string  `json:"text"`
// 	}
// }

// type EventMessageReceived struct {
// 	Type string `json:"type"`
// 	Body struct {
// 		ReplyTo *string `json:"replyTo,omitempty"`
// 		Channel *string `json:"channel,omitempty"`
// 		User    *string `json:"user,omitempty"`
// 		Text    string  `json:"text"`
// 	}
// }
