package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/ksaraf07/distributed-kv-store/store"
)

type Server struct {
	store     *store.Store
	followers []net.Conn
}

// New creates a Server backed by s, and dials every address in
// followerAddrs up front, keeping each connection open for replication.
func New(s *store.Store, followerAddrs []string) *Server {
	srv := &Server{store: s}

	for _, addr := range followerAddrs {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("could not connect to follower", addr, ":", err)
			continue
		}
		fmt.Println("connected to follower:", addr)
		srv.followers = append(srv.followers, conn)
	}

	return srv
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

// replicate forwards a command line to every connected follower.
func (srv *Server) replicate(command string) {
	for _, conn := range srv.followers {
		fmt.Fprintln(conn, command)
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
		if err := srv.store.Set(parts[1], parts[2]); err != nil {
			return "ERR " + err.Error()
		}
		srv.replicate(line)
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
		if err := srv.store.Delete(parts[1]); err != nil {
			return "ERR " + err.Error()
		}
		srv.replicate(line)
		return "OK"

	default:
		return "ERR unknown command"
	}
}
