package main

import (
	"fmt"
	"os"

	"github.com/ksaraf07/distributed-kv-store/server"
	"github.com/ksaraf07/distributed-kv-store/store"
)

func main() {
	s := store.New()
	srv := server.New(s)

	addr := ":9000"
	if err := srv.ListenAndServe(addr); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
