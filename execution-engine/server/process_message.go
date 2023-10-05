package server

import (
	"encoding/json"
	"log"

	"github.com/3ssalunke/leetcode-clone-exen/services"
	"github.com/streadway/amqp"
)

func (server *Server) ProcessMessage(msg amqp.Delivery) error {
	log.Printf("Received a message: %s", msg.Body)

	var payload *services.ExecutionPayload
	err := json.Unmarshal(msg.Body, payload)
	if err != nil {
		return err
	}

	_, err = services.ExecuteCode(payload)
	if err != nil {
		return err
	}

	return nil
}
