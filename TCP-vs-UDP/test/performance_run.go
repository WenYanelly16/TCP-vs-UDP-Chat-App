//Filename: ./test/performance_test.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/WenYanelly16/TCP-VS-UDP/pkg"
)

const (
	tcpAddr         = "localhost:8080"
	udpAddr         = "localhost:8088"
	clients         = 50
	msgs            = 100
	benchIterations = 1000
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--bench" {
		runBenchmarks()
		return
	}

	runPerformanceTest()
}

func runPerformanceTest() {
	fmt.Println("Running performance comparison...")

	tcpMetrics := testProtocol("tcp")
	udpMetrics := testProtocol("udp")

	printResults("TCP", tcpMetrics)
	printResults("UDP", udpMetrics)
}

func runBenchmarks() {
	fmt.Println("Running benchmarks...")

	fmt.Println("\nTCP Benchmark:")
	tcpDuration := benchmarkTCP()
	fmt.Printf("Completed %d iterations in %v\n", benchIterations, tcpDuration)
	fmt.Printf("Average latency: %v/op\n", tcpDuration/time.Duration(benchIterations))

	fmt.Println("\nUDP Benchmark:")
	udpDuration := benchmarkUDP()
	fmt.Printf("Completed %d iterations in %v\n", benchIterations, udpDuration)
	fmt.Printf("Average latency: %v/op\n", udpDuration/time.Duration(benchIterations))
}

func benchmarkTCP() time.Duration {
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		fmt.Println("TCP connection error:", err)
		return 0
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	start := time.Now()

	for i := 0; i < benchIterations; i++ {
		if _, err := conn.Write([]byte("test")); err != nil {
			fmt.Println("TCP write error:", err)
			break
		}
		if _, err := conn.Read(buf); err != nil {
			fmt.Println("TCP read error:", err)
			break
		}
	}
	return time.Since(start)
}

func benchmarkUDP() time.Duration {
	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		fmt.Println("UDP resolve error:", err)
		return 0
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("UDP connection error:", err)
		return 0
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	start := time.Now()

	for i := 0; i < benchIterations; i++ {
		if _, err := conn.Write([]byte("test")); err != nil {
			fmt.Println("UDP write error:", err)
			break
		}
		if _, _, err := conn.ReadFromUDP(buf); err != nil {
			fmt.Println("UDP read error:", err)
			break
		}
	}
	return time.Since(start)
}

func testProtocol(proto string) *pkg.Metrics {
	metrics := pkg.NewMetrics()
	var wg sync.WaitGroup
	wg.Add(clients)

	for i := 0; i < clients; i++ {
		go func(id int) {
			defer wg.Done()
			if proto == "tcp" {
				testTCPClient(id, metrics)
			} else {
				testUDPClient(id, metrics)
			}
		}(i)
	}
	wg.Wait()
	return metrics
}

func testTCPClient(id int, m *pkg.Metrics) {
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		m.RecordDrop()
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
		m.RecordDrop()
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		m.RecordDrop()
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

func printResults(proto string, m *pkg.Metrics) {
	fmt.Printf("\n=== %s Results ===\n", proto)
	fmt.Printf("Messages Sent: %d\n", msgs*clients)
	fmt.Printf("Messages Received: %d\n", m.MessageCount)
	fmt.Printf("Packet Loss: %.2f%%\n", float64(msgs*clients-m.MessageCount)/float64(msgs*clients)*100)
	fmt.Printf("Average Latency: %v\n", m.AverageLatency())
	fmt.Printf("Max Latency: %v\n", m.MaxLatency)
	fmt.Printf("Min Latency: %v\n", m.MinLatency)
	fmt.Printf("Dropped Packets: %d\n", m.DroppedPackets)
}
