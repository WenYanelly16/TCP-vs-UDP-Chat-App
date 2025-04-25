//Filename: ./test/integration_test.go
package main

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTCPConnection(t *testing.T) {
	// Start TCP server in a goroutine
	go func() {
		server := NewServer(":9090")
		server.Start()
	}()
	time.Sleep(100 * time.Millisecond) // Wait for server to start

	// Test connection
	conn, err := net.Dial("tcp", ":9090")
	require.NoError(t, err)
	require.NotNil(t, conn)
	conn.Close()
}

func TestUDPConnection(t *testing.T) {
	// Start UDP server in a goroutine
	go func() {
		server := NewServer(":9091")
		server.Start()
	}()
	time.Sleep(100 * time.Millisecond) // Wait for server to start

	// Test connection
	addr, err := net.ResolveUDPAddr("udp", ":9091")
	require.NoError(t, err)

	conn, err := net.DialUDP("udp", nil, addr)
	require.NoError(t, err)
	require.NotNil(t, conn)
	conn.Close()
}