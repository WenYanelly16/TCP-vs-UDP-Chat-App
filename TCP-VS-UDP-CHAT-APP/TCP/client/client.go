//Filename: ../TCP/client/client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Start() error {
	defer c.conn.Close()

	// Get name first
	fmt.Print("Enter your name: ")
	nameScanner := bufio.NewScanner(os.Stdin)
	nameScanner.Scan()
	name := nameScanner.Text()

	// Send name to server
	if _, err := fmt.Fprintln(c.conn, name); err != nil {
		return err
	}

	// Start message receiver
	go func() {
		scanner := bufio.NewScanner(c.conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Handle message input
	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		msg := inputScanner.Text()
		if msg == "/quit" {
			return nil
		}
		if _, err := fmt.Fprintln(c.conn, msg); err != nil {
			return err
		}
	}

	return nil
}