package mq

import (
	"log"

	"github.com/3ssalunke/leetcode-clone-exen/util"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(config util.Config) (*RabbitMQ, error) {
	rabbitMQAddr := config.RabbitMQHost
	conn, err := amqp.Dial(rabbitMQAddr)
	if err != nil {
		return nil, err
	}
	log.Println("rabbitmq connection established")
	return &RabbitMQ{Conn: conn}, nil
}

func (mq *RabbitMQ) CreateChannel() error {
	ch, err := mq.Conn.Channel()
	if err != nil {
		return err
	}

	log.Println("rabbitmq channel established")
	mq.Channel = ch
	return nil
}

func (mq *RabbitMQ) DeclareQueue(name string) error {
	_, err := mq.Channel.QueueDeclare(name, true, false, false, false, nil)
	return err
}

func (mq *RabbitMQ) Consume(queueName string) error {
	msgs, err := mq.Channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to setup a consumer - %v", err)
	}

	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)

		if err := msg.Ack(false); err != nil {
			log.Printf("Failed to acknowledge message: %v", err)
		}
	}

	return nil
}
