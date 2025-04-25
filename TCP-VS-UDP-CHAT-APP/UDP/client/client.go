//Filename: ../UDP/client/client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	
	"github.com/WenYanelly16/TCP-VS-UDP-CHAT-APP/pkg"

)

type Client struct {
	conn      *net.UDPConn
	serverAddr *net.UDPAddr
	done      chan struct{}
	messages  chan pkg.Message
}

func NewClient(addr string) *Client {
	return &Client{
		done:     make(chan struct{}),
		messages: make(chan pkg.Message, 100),
	}
}

func (c *Client) Start() error {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	c.serverAddr = serverAddr

	localAddr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		return err
	}
	c.conn = conn

	// Get client name
	fmt.Print("Enter your name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	if _, err := conn.Write([]byte("NAME:" + name)); err != nil {
		return err
	}

	go c.receiveMessages()
	go c.sendKeepAlives()

	for scanner.Scan() {
		text := scanner.Text()
		if text == "/quit" {
			break
		}
		if _, err := conn.Write([]byte(text)); err != nil {
			return err
		}
	}

	close(c.done)
	conn.Close()
	return nil
}

func (c *Client) receiveMessages() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-c.done:
			return
		default:
			n, _, err := c.conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			msg := string(buf[:n])
			if msg != "PING" {
				fmt.Println(msg)
			}
		}
	}
}

func (c *Client) sendKeepAlives() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.conn.Write([]byte("PING"))
		case <-c.done:
			return
		}
	}
}