//Filename: ../TCP/client/client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

type Client struct {
	conn     net.Conn
	mu       sync.Mutex
	done     chan struct{}
	messages chan pkg.Message
}

func NewClient(addr string) *Client {
	return &Client{
		done:     make(chan struct{}),
		messages: make(chan pkg.Message, 100),
	}
}

func (c *Client) Start() error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	c.conn = conn

	go c.receiveMessages()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "/quit" {
			break
		}
		if _, err := fmt.Fprintln(c.conn, text); err != nil {
			return err
		}
	}

	close(c.done)
	c.conn.Close()
	return nil
}

func (c *Client) receiveMessages() {
	scanner := bufio.NewScanner(c.conn)
	for {
		select {
		case <-c.done:
			return
		default:
			if scanner.Scan() {
				text := scanner.Text()
				fmt.Println(text)
			} else {
				return
			}
		}
	}
}