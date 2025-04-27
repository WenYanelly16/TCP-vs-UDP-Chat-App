//Filename: ../TCP/client/main.go

package main

import (
	"log"
)

func main() {
	client, err := NewClient("localhost:8080")
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	
	if err := client.Start(); err != nil {
		log.Fatal("Client error:", err)
	}
}