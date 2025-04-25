//Filename: /pkg/message.go
package pkg

import (
	"time"
)

type Message struct {
	From    string
	Content string
	Time    time.Time
}

func NewMessage(from, content string) Message {
	return Message{
		From:    from,
		Content: content,
		Time:    time.Now(),
	}
}

func (m Message) String() string {
	return m.Time.Format("15:04:05") + " " + m.From + ": " + m.Content
}