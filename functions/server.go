package functions

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Server manages the chat service, handling client connections, messages, and name registrations.
type Server struct {
	Listener        net.Listener         // TCP server listener
	Clients         map[net.Conn]string // Active clients 
	RegisteredNames map[string]bool      // Set of names in use
	MaxClients      int                  // Maximum number of clients
	MessageHistory  []string             // Stores all chat messages
	mutex           sync.Mutex           // Protects shared resources from concurrent access
	isShuttingDown  bool                  // Flag to indicate if server is shutting down
}

// NewServer - Initialize a new chat server.
func (s *Server) NewServer(port string, maxClients int) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s.Listener = listener
	s.Clients = make(map[net.Conn]string)
	s.RegisteredNames = make(map[string]bool)
	s.MaxClients = maxClients
	s.MessageHistory = []string{}
	return nil
}

// AllowConnection - Checks if a new connection can be accepted.
func (s *Server) AllowConnection(conn net.Conn) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.Clients) < s.MaxClients {
		return true
	}
	return false
}

// HandleConnection - Manages the connection.
func (s *Server) HandleConnection(conn net.Conn) {
	if !s.AllowConnection(conn) {
		fmt.Fprint(conn, "Chat room is full, please try again later...\n")
		conn.Close()
		return
	}

	fmt.Fprint(conn, GreetingMessage)

	var username string
	for {
		usernameInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Error reading username: %v", err)
			return
		}
		username = strings.TrimSpace(usernameInput)

		// Attempt to register the client with the provided username
		if err := s.registerClient(conn, username); err == nil {
			break // Registration successful, exit the loop
		} else {
			fmt.Fprint(conn, err.Error()+"\n[ENTER YOUR NAME]: ")
		}
	}

	s.startChat(conn)
	s.deregisterClient(conn)
}

// Close - Shutdown the server and notify all clients.
func (s *Server) Close() {
	s.isShuttingDown = true
	fmt.Println()
	fmt.Println(ColorRed + "Shutting down server..." + ColorReset)
	s.mutex.Lock()
	for conn := range s.Clients {
		fmt.Fprintf(conn, ColorRed + "\nServer is shutting down!" + ColorReset)
		conn.Close()
	}
	s.mutex.Unlock()
	s.Listener.Close()
	fmt.Println(ColorRed + "Server has been closed..." + ColorReset)
}

// startChat - Initiates the chat session.
func (s *Server) startChat(conn net.Conn) {
	s.sendMessageHistory(conn)
	message := Message(s, conn, "", ModeJoined)
	s.sendMessage(conn, message)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		formattedMessage := Message(s, conn, message, ModeSend)
		s.sendMessage(conn, formattedMessage)
		s.saveMessage(formattedMessage)
	}

	message = Message(s, conn, "", ModeLeft)
	s.sendMessage(conn, message)
}


// sendMessage - Broadcasts a message to all clients.
func (s *Server) sendMessage(conn net.Conn, message string) {
	timeStamp := time.Now().Format(TimeFormat)
	if message == "" {
		fmt.Fprintf(conn, UserPrompt, timeStamp, s.Clients[conn])
		return
	}
	// SendingMessage
	Message := fmt.Sprintf("%s%s\n%s", ColorRed, ColorReset, message)
	s.mutex.Lock()
	for connection := range s.Clients {
		if connection != conn {
			fmt.Fprint(connection, Message)
		}
		fmt.Fprintf(connection, UserPrompt, timeStamp, s.Clients[connection])
	}
	s.mutex.Unlock()
	
}

// sendMessageHistory - Sends all previous messages to the new client.
func (s *Server) sendMessageHistory(conn net.Conn) {
	s.mutex.Lock()
	for _, msg := range s.MessageHistory {
		fmt.Fprint(conn, msg)
	}
	s.mutex.Unlock()
}

// saveMessage - Appends a new message to the history.
func (s *Server) saveMessage(message string) {
	s.mutex.Lock()
	s.MessageHistory = append(s.MessageHistory, message)
	s.mutex.Unlock()
}

// registerClient - Adds a new client connection.
func (s *Server) registerClient(conn net.Conn, username string) error {
	if username == "" {
		return errors.New("Username cannot be empty")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.MaxClients > 0 && len(s.Clients) >= s.MaxClients {
		return fmt.Errorf("Chat room is full")
	}
	if s.RegisteredNames[username] {
		return fmt.Errorf("Username '%s' is already taken", username)
	}

	s.RegisteredNames[username] = true
	s.Clients[conn] = username
	log.Printf("Client connected: %v", conn.RemoteAddr())
	return nil
}

// deregisterClient - Removes a client connection.
func (s *Server) deregisterClient(conn net.Conn) {
	s.mutex.Lock()
	if s.isShuttingDown {
		// If server is shutting down, do not log disconnection
		delete(s.Clients, conn)
		return
	}
	username := s.Clients[conn]
	delete(s.RegisteredNames, username)
	delete(s.Clients, conn)
	log.Printf("Client disconnected: %v", conn.RemoteAddr())
	s.mutex.Unlock()
}
