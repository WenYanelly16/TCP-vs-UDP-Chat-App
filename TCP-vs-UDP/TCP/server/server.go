//Filename: ../TCP/server/server.go
package main

import (
	"net"
	"sync"

	"github.com/WenYanelly16/TCP-VS-UDP/pkg"
)

type Server struct {
	addr      string
	listener  net.Listener
	clients   map[net.Conn]*Client
	mu        sync.Mutex
	metrics   *pkg.Metrics
	broadcast chan pkg.Message
}

func NewServer(addr string) *Server {
	return &Server{
		addr:      addr,
		clients:   make(map[net.Conn]*Client),
		metrics:   pkg.NewMetrics(),
		broadcast: make(chan pkg.Message, 100),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener

	go s.handleBroadcasts()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && !ne.Temporary() {
				return err
			}
			continue
		}

		client := NewClient(conn)
		go s.handleConnection(client)
	}
}

func (s *Server) handleBroadcasts() {
	for msg := range s.broadcast {
		s.mu.Lock()
		for conn := range s.clients {
			if _, err := conn.Write([]byte(msg.String() + "\n")); err != nil {
				delete(s.clients, conn)
				conn.Close()
			}
		}
		s.mu.Unlock()
	}
}

func (s *Server) Stop() {
	close(s.broadcast)
	s.listener.Close()
}