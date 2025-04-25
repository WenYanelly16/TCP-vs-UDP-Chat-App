//Filename: ../UDP/server/server.go
package main

import (
	"net"
	"sync"
	"time"

	"github.com/WenYanelly16/TCP-VS-UDP-CHAT-APP/pkg"
)

type Server struct {
	addr      string
	conn      *net.UDPConn
	clients   map[string]*Client
	mu        sync.Mutex
	metrics   *pkg.Metrics
	broadcast chan pkg.Message
}

func NewServer(addr string) *Server {
	return &Server{
		addr:      addr,
		clients:   make(map[string]*Client),
		metrics:   pkg.NewMetrics(),
		broadcast: make(chan pkg.Message, 100),
	}
}

func (s *Server) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = conn

	go s.handleBroadcasts()
	go s.cleanupInactiveClients()

	buf := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		msg := string(buf[:n])
		go s.handleMessage(addr, msg)
	}
}

func (s *Server) handleBroadcasts() {
	for msg := range s.broadcast {
		s.mu.Lock()
		for _, client := range s.clients {
			if _, err := s.conn.WriteToUDP([]byte(msg.String()), client.addr); err != nil {
				delete(s.clients, client.addr.String())
			}
		}
		s.mu.Unlock()
	}
}

func (s *Server) cleanupInactiveClients() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for id, client := range s.clients {
			if time.Since(client.lastActive) > 5*time.Minute {
				s.broadcast <- pkg.NewMessage("Server", client.name+" timed out")
				delete(s.clients, id)
			}
		}
		s.mu.Unlock()
	}
}

func (s *Server) Stop() {
	close(s.broadcast)
	s.conn.Close()
}