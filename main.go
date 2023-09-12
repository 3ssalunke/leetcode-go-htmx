package main

import (
	"log"

	"github.com/3ssalunke/leetcode-clone/server"
)

func main() {
	log.Println("Let's Get Started!!")
	server := server.NewServer()
	log.Fatal(server.Start("127.0.0.1"))
}
