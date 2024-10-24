package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	joinMessage := fmt.Sprintf("JOIN:%s", name)
	conn.Write([]byte(joinMessage))

	go receiveMessages(conn)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "/leave" {
			leaveMessage := fmt.Sprintf("LEAVE:%s", name)
			conn.Write([]byte(leaveMessage))
			fmt.Println("You have left the chat")
			return
		}

		message := fmt.Sprintf("MSG:%s:%s", name, text)
		conn.Write([]byte(message))
	}
}

func receiveMessages(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		message := string(buffer[:n])
		fmt.Println(message)
	}
}
