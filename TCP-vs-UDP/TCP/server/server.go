//Filename: ../TCP/server/server.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
	"syscall" // Added this import

	"github.com/WenYanelly16/TCP-VS-UDP/pkg"
)

type Client struct {
	conn net.Conn
	name string
}

type Server struct {
	addr      string
	listener  net.Listener
	mu        sync.Mutex
	clients   map[net.Conn]*Client
	broadcast chan pkg.Message
	quit      chan struct{}
}

func NewServer(addr string) *Server {
	return &Server{
		addr:      addr,
		clients:   make(map[net.Conn]*Client),
		broadcast: make(chan pkg.Message, 100),
		quit:      make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.listener = ln

	//log.Printf("TCP Chat Server starting on %s", s.addr)

	// Broadcast loop
	go s.broadcastMessages()

	// Handle shutdown signals
	go s.handleSignals()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				log.Println("Accept error:", err)
				continue
			}
		}

		s.mu.Lock()
		s.clients[conn] = &Client{conn: conn}
		s.mu.Unlock()

		go s.handleConnection(conn)
	}
}

func (s *Server) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	s.Stop()
	os.Exit(0)
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
}

func (s *Server) broadcastMessages() {
	for {
		select {
		case msg := <-s.broadcast:
			s.mu.Lock()
			for conn := range s.clients {
				if _, err := conn.Write([]byte(msg.String() + "\n")); err != nil {
					log.Println("Broadcast error:", err)
				}
			}
			s.mu.Unlock()
		case <-s.quit:
			return
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		if client, ok := s.clients[conn]; ok {
			delete(s.clients, conn)
			s.broadcast <- pkg.NewMessage("Server", 
				time.Now().Format("15:04:05")+" "+client.name+" has left the chat")
		}
		s.mu.Unlock()
		conn.Close()
	}()

	client := s.clients[conn]
	reader := bufio.NewReader(conn)

	// Get client name
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	client.name = strings.TrimSpace(name)
	if client.name == "" {
		client.name = "Anonymous"
	}
	// Send welcome message
welcomeMsg := fmt.Sprintf(`
====================================
WELCOME TO THE CHAT, %s!
You are now connected to the server.
====================================
`, client.name)

if _, err := conn.Write([]byte(welcomeMsg)); err != nil {
    log.Printf("Error sending welcome to %s: %v", client.name, err)
    return
}
	// Broadcast join message
	s.broadcast <- pkg.NewMessage("Server", client.name+" has joined the chat")


	// Handle client messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		if message == "/quit" {
			break
		}

		s.broadcast <- pkg.NewMessage(client.name, message)
	}
}