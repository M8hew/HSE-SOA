package main

import (
	"flag"
	"log"
	"stat_service/internal/handlers"
)

func main() {
	livenessPort := flag.Int("liveliness_port", 8082, "Port number to listen on for liveness handler")
	grpcPort := flag.String("rpc_port", "8083", "Port number to listen on")
	flag.Parse()

	// liveness handler
	go func() {
		err := handlers.StartLivenessHandler(*livenessPort)
		if err != nil {
			log.Fatalf("Error running liveness handler: %v", err.Error())
		}
	}()

	// grpc handlers
	s := handlers.NewServer()
	if s == nil {
		log.Fatalln("Error creating new server")
	}
	log.Println("GRPC server successfully configured")
	log.Fatalf("GRPC server stopped: %v", s.ListenAndServe(*grpcPort))
}
