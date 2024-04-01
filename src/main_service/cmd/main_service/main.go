package main

import (
	"flag"
	"fmt"
	"main_service/api"
	"main_service/internal/handlers"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	configPath := flag.String("config_path", "", "Path to yaml configuration file")
	port := flag.String("port", "", "Port number to listen on")
	flag.Parse()

	e := echo.New()

	serverHandler, err := handlers.NewServerHandler(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting the server: %v", err)
		os.Exit(1)
	}
	api.RegisterHandlers(e, serverHandler)

	e.Start(":" + *port)
}
