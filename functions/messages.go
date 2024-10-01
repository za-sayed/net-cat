package functions

import (
	"fmt"
	"net"
	"time"
)

// Message - formats messages based on the given mode
func Message(server *Server, connection net.Conn, msg string, mode int) string {
	// Locking to safely access shared resources
	server.mutex.Lock()
	username := server.Clients[connection]
	server.mutex.Unlock()

	// Format the message based on the mode
	switch mode {
	case ModeSend:
		if msg == "\n" {
			return ""
		}
		timestamp := time.Now().Format(TimeFormat)
		msg = fmt.Sprintf(UserMessage, timestamp, username, msg)
	case ModeJoined:
		msg = fmt.Sprintf(ColorYellow + UserJoined + ColorReset, username)
	case ModeLeft:
		msg = fmt.Sprintf(ColorYellow + UserLeft + ColorReset, username)
	}
	return msg
}
