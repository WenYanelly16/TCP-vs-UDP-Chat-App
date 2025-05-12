
// Filename: ../TCP/server/server.go
// Filename: ../TCP/server/server.go
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
	"syscall"

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

	go s.broadcastMessages()
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
            msgBytes := []byte(msg.String() + "\n")
            deadConns := make([]net.Conn, 0)
            
            for conn := range s.clients {
                conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
                _, err := conn.Write(msgBytes)
                if err != nil {
                    if isConnectionError(err) {
                        deadConns = append(deadConns, conn)
                        continue
                    }
                    log.Printf("Broadcast error: %v", err)
                }
            }
            
            // Clean up dead connections
            for _, conn := range deadConns {
                if client, exists := s.clients[conn]; exists {
                    log.Printf("Client %s disconnected", client.name)
                    delete(s.clients, conn)
                    conn.Close()
                }
            }
            s.mu.Unlock()
            
        case <-s.quit:
            return
        }
    }
}

func isConnectionError(err error) bool {
    if opErr, ok := err.(*net.OpError); ok {
        return opErr.Err.Error() == "broken pipe" ||
               strings.Contains(opErr.Err.Error(), "connection reset") ||
               opErr.Timeout()
    }
    return false
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

	// Set initial deadline for handshake
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	client := s.clients[conn]
	reader := bufio.NewReader(conn)

	// Get client name
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading name: %v", err)
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

	// Reset deadline for normal operations
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Handle client messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Client %s timed out", client.name)
			} else {
				log.Printf("Client %s disconnected: %v", client.name, err)
			}
			break
		}

		message = strings.TrimSpace(message)
		if message == "/quit" {
			break
		}

		// Reset deadline after each successful read
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		s.broadcast <- pkg.NewMessage(client.name, message)
	}
}