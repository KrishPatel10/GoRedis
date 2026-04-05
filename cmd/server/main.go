package main

import (
	"fmt"

	"github.com/KrishPatel10/GoRedis.git/internal/aof"
	server "github.com/KrishPatel10/GoRedis.git/internal/server"
	store "github.com/KrishPatel10/GoRedis.git/internal/store"
)

func main() {
	addr := "0.0.0.0:6379"

	initServer(addr)
}

func initServer(addr string) {
	aof, err := aof.NewAOF("db.aof")

	if err != nil {
		fmt.Printf("Error while getting file pointer %s", err)
	}

	store := store.New(aof)

	s := server.NewServer(addr, store)

	s.Start()
}
