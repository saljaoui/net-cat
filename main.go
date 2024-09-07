package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	ListenAddr string
	ln         net.Listener
	clients    map[*Client]string
	users	map[string]bool
}

var allMessages string
var WelcomMessage = "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "


type Client struct {
	conn net.Conn
	Name string
}

func NewServer(ListenAddr string) *Server {
	return &Server{
		ListenAddr: ListenAddr,
		clients:    make(map[*Client]string),
		users: make(map[string]bool),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	s.acceptLoop()
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		client := &Client{conn: conn}
		go s.handleClient(client)
	}
}

func (s *Server) handleClient(client *Client) {
	client.conn.Write([]byte(WelcomMessage))
	name, err := bufio.NewReader(client.conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading client name:", err)
		return
	}

	client.Name = strings.TrimSpace(name)
	s.clients[client] = client.Name




	client.conn.Write([]byte(allMessages))

	s.broadcastMessage(fmt.Sprintf("%s has joined the chat\n", client.Name), client)

	r := bufio.NewReader(client.conn)

	for {

		msg, err := r.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from client %s: %v\n", client.Name, err)
			s.disconnectClient(client)
			return
		}
		s.broadcastMessage(formatMessage(msg, client.Name), client)
	}
}


func (s *Server) broadcastMessage(msg string, sender *Client) {
	saveMessages(msg)

	msg = "\n" + msg
	for client := range s.clients {
		if client != sender {
			_, err := client.conn.Write([]byte(msg))
			if err != nil {
				fmt.Printf("Error broadcasting to %s: %v\n", client.Name, err)
				client.conn.Close()
				delete(s.clients, client)
			}

		}
		client.conn.Write([]byte(formatMessage("", s.clients[client])))
	}
}

func (s *Server) disconnectClient(client *Client) {
    client.conn.Close()
    delete(s.clients, client)
    s.broadcastMessage(fmt.Sprintf("%s has left the chat\n", client.Name), nil)
}

func formatMessage(message string, name string) string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]: %s", formattedTime, name, message)
}

func saveMessages(msg string) {
	allMessages += msg
}

func main() {
	port := ":3000"
	fmt.Println("Server started on " + port)

	server := NewServer(port)
	err := server.Start()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
