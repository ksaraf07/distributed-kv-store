package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("could not connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("connected to kv-store. type commands like: SET key value")

	serverReader := bufio.NewScanner(conn)
	userInput := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !userInput.Scan() {
			break
		}
		line := userInput.Text()

		fmt.Fprintln(conn, line)

		if !serverReader.Scan() {
			fmt.Println("server closed the connection")
			break
		}
		fmt.Println(serverReader.Text())
	}

	if err := userInput.Err(); err != nil {
		fmt.Println("input error:", err)
	}
	if err := serverReader.Err(); err != nil {
		fmt.Println("connection error:", err)
	}
}
