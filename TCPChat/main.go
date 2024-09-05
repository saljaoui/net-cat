package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	port := "3000"

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on port " + port + "...\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			continue
		}
		fmt.Println(conn)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// s := make([]byte, 1000)
	fmt.Printf("Accepted connection from RemoteAddr %s\n", conn.RemoteAddr())
	fmt.Printf("Accepted connection from LocalAddr %s\n", conn.LocalAddr())
	data := bufio.NewReader(conn)
	conn.Write([]byte("Enter your name: "))
	name, _ := data.ReadString('\n')
	name = name[:len(name)-1]
	fmt.Println(name)
	
}
