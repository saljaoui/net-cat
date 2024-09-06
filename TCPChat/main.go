package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	Name string
}

var (
	clients = make(map[*Client]bool)

	Names = make(map[string]bool)

	clientsMux sync.Mutex
	messages   []string
)

func main() {
	port := ":8080"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("The following error occurred", err)
		return
	}
	fmt.Printf("Listening on the port %s\n",port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	client := &Client{conn: conn}

	conn.Write([]byte("[ENTER YOUR NAME]: "))
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		client.Name = scanner.Text()
	} else {
		return
	}

	clientsMux.Lock()

	clients[client] = true
	if Names[fmt.Sprint(client.Name)] {
		for {
			conn.Write([]byte("already here------------------\n"))
			conn.Write([]byte("[ENTER YOUR NAME]: "))
			scanner = bufio.NewScanner(conn)
			if scanner.Scan() {
				client.Name = scanner.Text()
			} else {
				return
			}
			if !Names[fmt.Sprint(client.Name)] {
				Names[fmt.Sprint(client.Name)] = true
				break
			}
		}
	} else {
		Names[fmt.Sprint(client.Name)] = true
	}

	for _, msg := range messages {
		conn.Write([]byte(msg))
	}

	clientsMux.Unlock()



	broadcastMessage(fmt.Sprintf("\n%s has joined our chat...\n", client.Name), client)

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	conn.Write([]byte(fmt.Sprintf("[%s][%s]: ", formattedTime, client.Name)))

	for scanner.Scan() {
		
		currentTime = time.Now()
		formattedTime = currentTime.Format("2006-01-02 15:04:05")
		
		conn.Write([]byte(fmt.Sprintf("[%s][%s]: ", formattedTime, client.Name)))

		message := scanner.Text()
		broadcastMessage(fmt.Sprintf("[%s][%s]: %s\n", formattedTime, client.Name, message), client)
	}

	clientsMux.Lock()
	delete(clients, client)
	clientsMux.Unlock()

	broadcastMessage(fmt.Sprintf("%s has left the chat...\n", client.Name), client)
	delete(Names, client.Name)

}


func broadcastMessage(message string, sender *Client) {
	fmt.Print(message)
	
	messages = append(messages, message)

	clientsMux.Lock()
	defer clientsMux.Unlock()

	for client := range clients {
		if client != sender { 
			_, err := client.conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("Error broadcasting to %s: %v\n", client.Name, err)
				client.conn.Close()
				delete(clients, client)
				
			}
		}
	}
}
