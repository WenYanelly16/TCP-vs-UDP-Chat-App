//Filename: ./test/performance_test.go
package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/WenYanelly16/TCP-VS-UDP/pkg"

)

const (
	tcpAddr = "localhost:8080"
	udpAddr = "localhost:8081"
	clients = 50
	msgs    = 100
)

func main() {
	fmt.Println("Running performance comparison...")

	tcpMetrics := pkg.NewMetrics()
	udpMetrics := pkg.NewMetrics()

	var wg sync.WaitGroup
	wg.Add(clients * 2)

	// TCP Test
	start := time.Now()
	for i := 0; i < clients; i++ {
		go func(id int) {
			defer wg.Done()
			testTCPClient(id, tcpMetrics)
		}(i)
	}

	// UDP Test
	for i := 0; i < clients; i++ {
		go func(id int) {
			defer wg.Done()
			testUDPClient(id, udpMetrics)
		}(i)
	}

	wg.Wait()

	tcpDuration := time.Since(start)
	printResults("TCP", tcpMetrics, tcpDuration)
	printResults("UDP", udpMetrics, tcpDuration)
}

func testTCPClient(id int, m *pkg.Metrics) {
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "NAME:client%d\n", id)
	reader := bufio.NewReader(conn)

	for i := 0; i < msgs; i++ {
		msg := fmt.Sprintf("msg%d", i)
		start := time.Now()
		fmt.Fprintln(conn, msg)

		if _, err := reader.ReadString('\n'); err != nil {
			m.RecordDrop()
			continue
		}
		m.Record(time.Since(start))
	}
}

func testUDPClient(id int, m *pkg.Metrics) {
	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "NAME:client%d\n", id)
	buf := make([]byte, 1024)

	for i := 0; i < msgs; i++ {
		msg := fmt.Sprintf("msg%d", i)
		start := time.Now()
		conn.Write([]byte(msg))

		if _, _, err := conn.ReadFromUDP(buf); err != nil {
			m.RecordDrop()
			continue
		}
		m.Record(time.Since(start))
	}
}

func printResults(proto string, m *pkg.Metrics, duration time.Duration) {
	fmt.Printf("\n=== %s Results ===\n", proto)
	fmt.Printf("Test Duration: %v\n", duration)
	fmt.Printf("Messages Sent: %d\n", msgs*clients)
	fmt.Printf("Messages Received: %d\n", m.MessageCount)
	fmt.Printf("Packet Loss: %.2f%%\n", float64(msgs*clients-m.MessageCount)/float64(msgs*clients)*100)
	fmt.Printf("Average Latency: %v\n", m.AverageLatency())
	fmt.Printf("Max Latency: %v\n", m.MaxLatency)
	fmt.Printf("Min Latency: %v\n", m.MinLatency)
}