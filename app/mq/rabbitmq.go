package mq

import (
	"log"

	"github.com/3ssalunke/leetcode-clone-app/util"
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

func (mq *RabbitMQ) DeclareQueue(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) error {
	_, err := mq.Channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	return err
}

func (mq *RabbitMQ) PublishMessage(queueName string, message string) error {
	err := mq.Channel.Publish("", queueName, false, false, amqp.Publishing{ContentType: "application/json", Body: []byte(message)})
	return err
}
