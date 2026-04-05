package server

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/KrishPatel10/GoRedis.git/internal/resp"
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

	parser := resp.NewParser(conn)
	writer := resp.NewWriter(conn)

	for {
		val, err := parser.Parse()

		if err != nil {
			if err == io.EOF {
				fmt.Print("Empty lines")
				break
			}
			fmt.Printf("Connection Error: %s", err)
		}

		if val.Typ != "array" || len(val.Array) == 0 {
			continue
		}

		command := strings.ToUpper(val.Array[0].Value)

		switch command {
		case "PING":
			writer.WriteSimpleString("PONG")

		case "SET":
			if len(val.Array) < 3 {
				continue
			}

			key := val.Array[1].Value
			val := val.Array[2].Value

			cache.SetWithoutExpiry(key, val)

			writer.WriteSimpleString("OK")

		case "GET":
			if len(val.Array) != 2 {
				continue
			}

			key := val.Array[1].Value

			res, exists := cache.Get(key)

			if !exists {
				writer.WriteNull()
			} else {
				writer.WriteBulkString(res) // Clean, exact data return
			}
		default:
			writer.WriteSimpleString("No Command found matching")
		}
	}
}
