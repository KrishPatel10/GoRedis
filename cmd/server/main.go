package main

import (
	server "github.com/KrishPatel10/GoRedis.git/internal/server"
	store "github.com/KrishPatel10/GoRedis.git/internal/store"
)

func main() {
	addr := "localhost:6379"

	initServer(addr)
}

func initServer(addr string) {
	store := store.New()

	s := server.NewServer(addr, store)

	s.Start()
}
