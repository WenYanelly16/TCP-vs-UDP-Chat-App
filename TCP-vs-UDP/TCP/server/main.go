//Filename: ../TCP/server/main.go
package main

import (
	"log"
)

func main() {
	server := NewServer(":8080")
	log.Println("TCP Chat Server starting on :8080")
	log.Fatal(server.Start())
}