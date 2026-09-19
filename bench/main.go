package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	addr := flag.String("addr", "localhost:9000", "node to benchmark")
	clients := flag.Int("clients", 10, "number of concurrent simulated clients")
	requests := flag.Int("requests", 100, "requests per client")
	flag.Parse()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalLatency time.Duration
	var totalRequests int

	start := time.Now()

	for c := 0; c < *clients; c++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", *addr)
			if err != nil {
				fmt.Println("client", clientID, "could not connect:", err)
				return
			}
			defer conn.Close()

			reader := bufio.NewScanner(conn)

			for i := 0; i < *requests; i++ {
				reqStart := time.Now()

				key := fmt.Sprintf("bench-%d-%d", clientID, i)
				fmt.Fprintf(conn, "SET %s value\n", key)
				reader.Scan() // wait for response

				latency := time.Since(reqStart)

				mu.Lock()
				totalLatency += latency
				totalRequests++
				mu.Unlock()
			}
			if err := reader.Err(); err != nil {
				fmt.Println("client", clientID, "read error:", err)
			}
		}(c)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("=== Benchmark Results ===")
	fmt.Println("Target:", *addr)
	fmt.Println("Concurrent clients:", *clients)
	fmt.Println("Requests per client:", *requests)
	fmt.Println("Total requests:", totalRequests)
	fmt.Println("Total time:", elapsed)
	fmt.Printf("Throughput: %.0f requests/sec\n", float64(totalRequests)/elapsed.Seconds())
	fmt.Printf("Average latency: %v\n", totalLatency/time.Duration(totalRequests))
}
