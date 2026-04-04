package server

import (
	"fmt"
	"net"

	store "github.com/KrishPatel10/GoRedis.git/internal/store"
)

type Server struct {
	addr  string
	store *store.MemoryStore
}

func NewServer(addr string, store *store.MemoryStore) *Server {
	return &Server{
		addr:  addr,
		store: store,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("Listening to port %s", s.addr)

	// infinite loop to listen request
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Failed to accept connection: %v", err)
			continue // don't break server for one connection break
		}

		go handleConnection(conn, s.store)
	}
}

func handleConnection(conn net.Conn, cache *store.MemoryStore) {
	defer conn.Close()

	fmt.Print("Got Connection\n")
}
