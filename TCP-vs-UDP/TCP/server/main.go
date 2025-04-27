//Filename: ../TCP/server/main.go
package main

import "log"

func main() {
	srv := NewServer(":8080")
	log.Printf("TCP Chat Server starting on %s", srv.addr)  
	log.Fatal(srv.Start())
}