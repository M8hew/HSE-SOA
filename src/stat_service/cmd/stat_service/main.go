package main

import (
	"flag"
	"log"
	"stat_service/internal/handlers"
)

func main() {
	livenessPort := flag.Int("port", 8082, "Port number to listen on for liveness handler")
	flag.Parse()

	err := handlers.StartLivenessHandler(*livenessPort)
	if err != nil {
		log.Fatalf("Error running liveness handler: %v", err.Error())
	}
}
