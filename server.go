package main

import (
	"fmt"
	"net"
	"strings"
)

var clients = make(map[string]*net.UDPAddr)

func main() {
	address, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Server running on port 8080...")

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		message := string(buffer[:n])
		parts := strings.SplitN(message, ":", 3) 

		if len(parts) >= 2 {
			command := strings.TrimSpace(parts[0])

			switch command {
			case "JOIN":
				name := strings.TrimSpace(parts[1])
				clients[name] = clientAddr
				broadcastMessage(fmt.Sprintf("%s has joined the chat", name), conn, clientAddr)
			case "LEAVE":
				name := strings.TrimSpace(parts[1])
				delete(clients, name)
				broadcastMessage(fmt.Sprintf("%s has left the chat", name), conn, clientAddr)
			case "MSG":
				if len(parts) == 3 {
					sender := strings.TrimSpace(parts[1])
					content := strings.TrimSpace(parts[2])
					broadcastMessage(fmt.Sprintf("%s: %s", sender, content), conn, clientAddr, sender)
				}
			}
		}
	}
}

func broadcastMessage(message string, conn *net.UDPConn, senderAddr *net.UDPAddr, senderName ...string) {
	for name, addr := range clients {
		if len(senderName) > 0 && name == senderName[0] {
			continue
		}
		_, err := conn.WriteToUDP([]byte(message), addr)
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", name, err)
		}
	}
	fmt.Println("Broadcasted:", message)
}
