package main

import (
	"content_service/internal/server"
	"flag"
	"fmt"
	"log"
)

func main() {
	port := flag.String("port", "", "Port number to listen on")
	flag.Parse()

	s := server.NewServer()
	if s == nil {
		log.Fatalln("Error creating new server")
	}
	fmt.Println("Server successfully configured")
	log.Fatalf("Server stopped: %v", s.ListenAndServe(*port))
}
