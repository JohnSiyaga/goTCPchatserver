package main

import (
	"log"
	"net"
)

// TCP Server
func main() {
	// Initialize new server
	s := newServer()
	go s.run()

	// Create listener on tcp channel 8888
	listener, err := net.Listen("tcp", ":8888")

	// Check for errors in server start
	if err != nil {
		log.Fatalf("Unable to start server: %s", err.Error())
	}

	// Close listener
	defer listener.Close()
	log.Printf("Started server on :8888")

	// Loop for incoming client connections
	for {
		conn, err := listener.Accept()
		// Check for errors in connection
		if err != nil {
			log.Printf("Unable to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
