package main

import (
	"log"

	"github.com/3ssalunke/leetcode-clone-app/db"
	"github.com/3ssalunke/leetcode-clone-app/server"
)

func main() {
	log.Println("starting app!!")
	server := server.NewServer()
	defer db.CloseMongoDbConnection(server.Db.Client())
	defer server.Mq.Conn.Close()
	log.Fatal(server.Start("127.0.0.1"))
}
