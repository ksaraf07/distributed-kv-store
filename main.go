package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ksaraf07/distributed-kv-store/server"
	"github.com/ksaraf07/distributed-kv-store/store"
)

func main() {
	port := flag.String("port", "9000", "port to listen on")
	logPath := flag.String("log", "kv.log", "path to the write-ahead log file")
	followers := flag.String("followers", "", "comma-separated list of follower addresses")
	flag.Parse()

	var followerAddrs []string
	if *followers != "" {
		followerAddrs = strings.Split(*followers, ",")
	}

	s, err := store.New(*logPath)
	if err != nil {
		fmt.Println("failed to open store:", err)
		os.Exit(1)
	}

	srv := server.New(s, followerAddrs)

	addr := ":" + *port
	if err := srv.ListenAndServe(addr); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
}
