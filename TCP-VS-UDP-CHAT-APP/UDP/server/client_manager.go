//Filename: ../UDP/server/client_manager.go
package main

import (
	"net"
	"strings"
	"time"

    "github.com/WenYanelly16/TCP-VS-UDP-CHAT-APP/pkg"
)

type Client struct {
	addr       *net.UDPAddr
	name       string
	lastActive time.Time
}

func (s *Server) handleMessage(addr *net.UDPAddr, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clientID := addr.String()
	client, exists := s.clients[clientID]

	if !exists {
		if strings.HasPrefix(msg, "NAME:") {
			name := strings.TrimPrefix(msg, "NAME:")
			s.clients[clientID] = &Client{
				addr:       addr,
				name:       name,
				lastActive: time.Now(),
			}
			s.broadcast <- pkg.NewMessage("Server", name+" has joined the chat")
		}
		return
	}

	client.lastActive = time.Now()

	if msg == "/quit" {
		s.broadcast <- pkg.NewMessage("Server", client.name+" has left the chat")
		delete(s.clients, clientID)
		return
	}

	if msg != "PING" {
		s.broadcast <- pkg.NewMessage(client.name, msg)
	}
}