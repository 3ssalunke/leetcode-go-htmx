package server

import (
	"log"

	"github.com/3ssalunke/leetcode-clone-exen/mq"
	"github.com/3ssalunke/leetcode-clone-exen/util"
)

type Server struct {
	config util.Config
	Mq     *mq.RabbitMQ
}

func NewServer() *Server {
	server := Server{}

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load enviroment variables - %v", err)
	}

	rabbitmq, err := mq.NewRabbitMQ(config)
	if err != nil {
		log.Fatalf("failed to make connection to rabbitmq server - %v", err)
	}

	err = rabbitmq.CreateChannel()
	if err != nil {
		log.Fatalf("failed to create channel to rabbitmq server - %v", err)
	}

	err = rabbitmq.DeclareQueue(config.RabbitMQQueueName)
	if err != nil {
		log.Fatalf("failed to declare a rabbitmq queue - %v", err)
	}

	server.Mq = rabbitmq
	server.config = config

	return &server
}

func (server *Server) StartExecutionEngine() error {
	return server.Mq.Consume(server.config.RabbitMQQueueName)
}
