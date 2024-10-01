package main

import (
	"fmt"
	"log"
	"net-cat/functions"
	"os"
	"os/signal"
	"syscall"
)

// Default settings
const (
	Port           = ":8989" // Default port for the server
	MaxConnections = 10      // Maximum number of clients
)

func main() {
	// Initialize the server
	port, err := getPort()
	if err != nil {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}

	server := &functions.Server{}
	if err := server.NewServer(port, MaxConnections); err != nil {
		log.Fatalf("ERROR -> main: %v\n", err)
		os.Exit(0)
	}

	fmt.Printf("Listening on port %v\n", port)

	// Set up signal handling for graceful shutdown
	closeSignal := make(chan os.Signal, 1)
	go handleCloseSignal(server, closeSignal)

	// Accept incoming connections
	acceptConnections(server)
	<-closeSignal
}

// getPortFromArgs retrieves the port from command line arguments
func getPort() (string, error) {
	if len(os.Args) < 2 {
		return Port, nil // Return default port if none provided
	} else if len(os.Args) > 2 {
		return "", fmt.Errorf("User inputs more than 1 argument")
	}
	return ":" + os.Args[1], nil
}

// handleCloseSignal listens for termination signals to gracefully shut down the server
func handleCloseSignal(server *functions.Server, closeSignal chan os.Signal) {
	signal.Notify(closeSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-closeSignal
	server.Close() // Close the server safely
	os.Exit(0)     // Exit the program
}

// acceptConnections listens for incoming connections and handles them
func acceptConnections(server *functions.Server) {
    for {
        conn, err := server.Listener.Accept()
        if err != nil {
            break
        }
        go server.HandleConnection(conn) // Handle connection in a new goroutine
    }
}

