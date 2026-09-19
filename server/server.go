package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ksaraf07/distributed-kv-store/store"
)

type Server struct {
	store     *store.Store
	followers []net.Conn

	mu            sync.Mutex
	isLeader      bool
	lastHeartbeat time.Time
}

// New creates a Server backed by s. If followerAddrs is non-empty, this
// node is treated as the leader: it dials each follower up front and
// begins sending periodic heartbeats. Every node also runs a monitor
// that watches for missed heartbeats and self-promotes if the leader
// appears to have died.
func New(s *store.Store, followerAddrs []string) *Server {
	srv := &Server{
		store:    s,
		isLeader: len(followerAddrs) > 0,
	}

	for _, addr := range followerAddrs {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("could not connect to follower", addr, ":", err)
			continue
		}
		fmt.Println("connected to follower:", addr)
		srv.followers = append(srv.followers, conn)
	}

	if srv.isLeader {
		go srv.sendHeartbeats()
	}
	go srv.monitorLeader()

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

// sendHeartbeats runs only on the leader, pinging every follower once
// per second over the same connections used for replication.
func (srv *Server) sendHeartbeats() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		srv.replicate("PING")
	}
}

// recordHeartbeat is called whenever this node receives a PING,
// marking "the leader is still alive as of right now".
func (srv *Server) recordHeartbeat() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.lastHeartbeat = time.Now()
}

// monitorLeader runs on every node. If this node is already the leader,
// it does nothing. Otherwise, it checks every 500ms whether it's been
// too long since the last heartbeat — if so, it promotes itself.
func (srv *Server) monitorLeader() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		srv.mu.Lock()
		isLeader := srv.isLeader
		last := srv.lastHeartbeat
		srv.mu.Unlock()

		if isLeader {
			continue
		}
		if last.IsZero() {
			continue // never heard from a leader yet
		}
		if time.Since(last) > 3*time.Second {
			srv.mu.Lock()
			srv.isLeader = true
			srv.mu.Unlock()
			fmt.Println("leader appears down — promoting self to leader")
		}
	}
}

func (srv *Server) handleCommand(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "ERR empty command"
	}

	switch strings.ToUpper(parts[0]) {
	case "PING":
		srv.recordHeartbeat()
		return "PONG"

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
