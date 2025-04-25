//Filename: ../TCP/server/main.go
package main

import (
	"log"
	//"github.com/WenYanelly16/TCP-VS-UDP-CHAT-APP/TCP/pkg"
)

func main() {
	srv := NewServer(":8080")
	log.Println("TCP Chat Server starting on :8080")
	log.Fatal(srv.Start())
}