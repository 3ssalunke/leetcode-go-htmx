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

func (mq *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return mq.Channel.Consume(queueName, "", false, false, false, false, nil)

}
