//Filename: ../UDP/server/main.go
package main

import (
	"log"
	"net"
	
)

func main() {
    log.Println("UDP Server starting on :8088")
    pc, err := net.ListenPacket("udp", ":8088")
    if err != nil {
        log.Fatal(err)
    }
    defer pc.Close()

    for {
        buf := make([]byte, 1024)
        n, addr, err := pc.ReadFrom(buf)
        if err != nil {
            log.Printf("Error reading: %v", err)
            continue
        }
        log.Printf("Received %d bytes from %s: %s", n, addr, string(buf[:n]))
        
        // Echo back (optional)
        pc.WriteTo(buf[:n], addr)
    }
}