package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	ListenAddr string
	ln         net.Listener
	quitch     chan struct{}
}

func NewServer(ListenAddr string) *Server {
	return &Server{
		ListenAddr: ListenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	return nil
}

func (s *Server) acceptLoop() {
	
	for {
		con, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go s.readLoop(con)
	}

}

func (s *Server) readLoop(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)
	
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Fatal("error from reading")
		}
	}
}

func main() {
	server := NewServer(":3000")
	server.Start()
}
