//Filename: ../UDP/client/main.go
package main

import (
	"log"
	"os"
)

func main() {
	client := NewClient("localhost:8081")
	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}