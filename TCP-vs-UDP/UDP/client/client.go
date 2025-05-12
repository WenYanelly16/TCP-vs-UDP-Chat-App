//Filename: ../UDP/client/client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type UDPClient struct {
	serverAddr *net.UDPAddr
	conn       *net.UDPConn
	name       string
}

func NewUDPClient(addr string) (*UDPClient, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	// Use nil for local address to let system choose
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &UDPClient{
		serverAddr: udpAddr,
		conn:       conn,
	}, nil
}

func (c *UDPClient) Start() error {
	defer c.conn.Close()

	// Get user name
	fmt.Print("Enter your name: ")
	name, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}
	c.name = strings.TrimSpace(name)

	// Send name to server
	if _, err := c.conn.Write([]byte(c.name + "\n")); err != nil {
		return err
	}

	// Handle signals for clean shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		c.conn.Write([]byte("/quit\n"))
		c.conn.Close()
		os.Exit(0)
	}()

	// Message receiver
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := c.conn.Read(buf)
			if err != nil {
				fmt.Println("\nDisconnected from server")
				os.Exit(0)
			}
			fmt.Print(string(buf[:n]))
		}
	}()

	// Message sender
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if _, err := c.conn.Write([]byte(text + "\n")); err != nil {
			return err
		}
		if text == "/quit" {
			return nil
		}
	}

	return scanner.Err()
}