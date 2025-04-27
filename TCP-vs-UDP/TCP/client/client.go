//Filename: /TCP/server/client.go
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

type Client struct {
	conn net.Conn
	name string
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Start() error {
	defer c.conn.Close()

	// Get user name
	fmt.Print("Enter your name: ")
	name, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return fmt.Errorf("name input error: %w", err)
	}
	c.name = strings.TrimSpace(name)

	// Send name to server
	if _, err := c.conn.Write([]byte(c.name + "\n")); err != nil {
		return fmt.Errorf("send name error: %w", err)
	}

	// Setup signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		c.conn.Close()
		os.Exit(0)
	}()

	// Message receiver
	go func() {
		reader := bufio.NewReader(c.conn)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nDisconnected from server")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	// Message sender
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "/quit" {
			return nil
		}
		if _, err := c.conn.Write([]byte(text + "\n")); err != nil {
			return fmt.Errorf("send message error: %w", err)
		}
	}

	return scanner.Err()
}