//Filename: ../UDP/server/server.go
package main

import (
    "fmt"
    "log"
    "net"
    "strings"
    "sync"

    "github.com/WenYanelly16/TCP-VS-UDP/pkg"
)

type UDPClient struct {
    addr *net.UDPAddr
    name string
}

type UDPServer struct {
    addr      string
    conn      *net.UDPConn
    mu        sync.Mutex
    clients   map[string]*UDPClient // Keyed by addr.String()
    broadcast chan pkg.Message
    quit      chan struct{}
}

func NewUDPServer(addr string) *UDPServer {
    return &UDPServer{
        addr:      addr,
        clients:   make(map[string]*UDPClient),
        broadcast: make(chan pkg.Message, 100),
        quit:      make(chan struct{}),
    }
}

func (s *UDPServer) Start() error {
    udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
    if err != nil {
        return err
    }

    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return err
    }
    defer conn.Close()
    s.conn = conn

    go s.broadcastMessages()
    go s.receiveMessages()

    <-s.quit
    return nil
}

func (s *UDPServer) broadcastMessages() {
    buf := make([]byte, 1024)
    for msg := range s.broadcast {
        s.mu.Lock()
        text := msg.String() + "\n"
        copy(buf, []byte(text))
        
        for _, client := range s.clients {
            if _, err := s.conn.WriteToUDP(buf[:len(text)], client.addr); err != nil {
                log.Printf("Broadcast error to %s: %v", client.addr.String(), err)
            }
        }
        s.mu.Unlock()
    }
}

func (s *UDPServer) receiveMessages() {
    buf := make([]byte, 1024)
    for {
        select {
        case <-s.quit:
            return
        default:
            n, addr, err := s.conn.ReadFromUDP(buf)
            if err != nil {
                log.Printf("Read error: %v", err)
                continue
            }

            msg := strings.TrimSpace(string(buf[:n]))
            s.handleMessage(addr, msg)
        }
    }
}

func (s *UDPServer) handleMessage(addr *net.UDPAddr, msg string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    clientKey := addr.String()
    client, exists := s.clients[clientKey]

    if !exists {
        // New client joining
        name := strings.TrimSpace(msg)
        if name == "" {
            name = "Anonymous"
        }

        client = &UDPClient{
            addr: addr,
            name: name,
        }
        s.clients[clientKey] = client

        // Send welcome message
        welcome := fmt.Sprintf("WELCOME %s! There are %d users online\n", client.name, len(s.clients))
        s.conn.WriteToUDP([]byte(welcome), addr)

        // Broadcast join message
        s.broadcast <- pkg.NewMessage("Server", 
            fmt.Sprintf("%s has joined the chat", client.name))
        return
    }

    if msg == "/quit" {
        delete(s.clients, clientKey)
        s.broadcast <- pkg.NewMessage("Server",
            fmt.Sprintf("%s has left the chat", client.name))
        return
    }

    // Regular message
    s.broadcast <- pkg.NewMessage(client.name, msg)
}

func (s *UDPServer) Stop() {
    close(s.quit)
    s.conn.Close()
}