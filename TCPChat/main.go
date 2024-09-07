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
	clients       = make(map[*Client]bool)
	WelcomMessage = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "
	Names         = make(map[string]bool)

	clientsMux sync.Mutex
	messages []string
)

func main() {
	port := ":8080"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("The following error occurred", err)
		return
	}
	fmt.Printf("Listening on the port %s\n", port)

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

	conn.Write([]byte(WelcomMessage))

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
			conn.Write([]byte("Name already taken. Please choose another name.\n"))
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

	defer conn.Close()
	for _, msg := range messages {
		conn.Write([]byte(msg))
	}

	clientsMux.Unlock()

	broadcastMessage(fmt.Sprintf("\n%s has joined our chat...", client.Name), client)

	conn.Write([]byte(formatMessage("", client.Name)))
	for scanner.Scan() {

		conn.Write([]byte(formatMessage("", client.Name)))

		message := scanner.Text()
		
		broadcastMessage(formatMessage(message, client.Name), client)

	}

	clientsMux.Lock()
	delete(clients, client)
	clientsMux.Unlock()
	delete(Names, client.Name)

	broadcastMessage(fmt.Sprintf("%s has left the chat...", client.Name), client)
}

func formatMessage(message string, name string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]: %s", formattedTime, name, message)
}

func broadcastMessage(message string, sender *Client) {
	message = message+"\n"
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
