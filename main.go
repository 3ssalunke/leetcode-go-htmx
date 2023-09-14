package main

import (
	"log"

	"github.com/3ssalunke/leetcode-clone/db"
	"github.com/3ssalunke/leetcode-clone/server"
)

func main() {
	log.Println("Let's Get Started!!")
	server := server.NewServer()
	defer db.CloseMongoDbConnection(server.Db.Client())
	log.Fatal(server.Start("127.0.0.1"))
}
