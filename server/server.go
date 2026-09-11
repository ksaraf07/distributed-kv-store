package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/ksaraf07/distributed-kv-store/store"
)

type Server struct {
	store *store.Store
}

func New(s *store.Store) *Server {
	return &Server{store: s}
}

func (srv *Server) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("kv-store listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go srv.handleConn(conn)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		response := srv.handleCommand(line)
		fmt.Fprintln(conn, response)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("connection error:", err)
	}
}

func (srv *Server) handleCommand(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "ERR empty command"
	}

	switch strings.ToUpper(parts[0]) {
	case "SET":
		if len(parts) != 3 {
			return "ERR usage: SET key value"
		}
		srv.store.Set(parts[1], parts[2])
		return "OK"

	case "GET":
		if len(parts) != 2 {
			return "ERR usage: GET key"
		}
		value, ok := srv.store.Get(parts[1])
		if !ok {
			return "(nil)"
		}
		return value

	case "DEL":
		if len(parts) != 2 {
			return "ERR usage: DEL key"
		}
		srv.store.Delete(parts[1])
		return "OK"

	default:
		return "ERR unknown command"
	}
}
