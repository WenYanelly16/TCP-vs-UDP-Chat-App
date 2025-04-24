//Filename: ../TCP/server/client_handler.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/yourusername/go-chat-comparison/pkg"
)

type Client struct {
	conn   net.Conn
	name   string
	server *Server
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (s *Server) handleConnection(client *Client) {
	defer func() {
		s.mu.Lock()
		delete(s.clients, client.conn)
		s.mu.Unlock()
		client.conn.Close()
	}()

	s.mu.Lock()
	s.clients[client.conn] = client
	s.mu.Unlock()

	// Get client name
	fmt.Fprint(client.conn, "Enter your name: ")
	name, err := bufio.NewReader(client.conn).ReadString('\n')
	if err != nil {
		return
	}
	client.name = strings.TrimSpace(name)

	// Notify others
	s.broadcast <- pkg.NewMessage("Server", client.name+" has joined the chat")

	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "/quit" {
			break
		}
		s.broadcast <- pkg.NewMessage(client.name, text)
	}

	s.broadcast <- pkg.NewMessage("Server", client.name+" has left the chat")
}