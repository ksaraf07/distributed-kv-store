package main

import (
	"fmt"

	"github.com/ksaraf07/distributed-kv-store/store"
)

func main() {
	s := store.New()

	s.Set("name", "kush")
	value, ok := s.Get("name")
	fmt.Println("get name:", value, ok)

	s.Delete("name")
	value, ok = s.Get("name")
	fmt.Println("get name after delete:", value, ok)
}
