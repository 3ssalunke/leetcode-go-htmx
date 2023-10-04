package main

import (
	"log"

	"github.com/3ssalunke/leetcode-clone-exen/server"
)

func main() {
	log.Println("starting execution engine!!")
	server := server.NewServer()
	defer server.Mq.Conn.Close()
	log.Fatal(server.StartExecutionEngine())
}
