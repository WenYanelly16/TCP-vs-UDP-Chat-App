//Filename: ../UDP/server/main.go
package main

import (
	"log"
)

func main() {
	server := NewServer(":8081")
	log.Println("UDP Chat Server starting on :8081")
	log.Fatal(server.Start())
}