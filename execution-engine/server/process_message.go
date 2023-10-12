package server

import (
	"encoding/json"
	"log"

	"github.com/3ssalunke/leetcode-clone-exen/services"
	"github.com/streadway/amqp"
)

func (server *Server) ProcessMessage(msg amqp.Delivery) {
	var payload services.ProblemDetails
	err := json.Unmarshal(msg.Body, &payload)
	if err != nil {
		log.Printf("failed to parse message json %v", err)
		return
	}

	if err := server.Redis.SetValue(payload.ExecutionId, "processing"); err != nil {
		log.Printf("execution id: %s, failed to set value in redis cache - %v", payload.ExecutionId, err)
		return
	}

	result, err := services.ExecuteCode(&payload)
	if err != nil {
		if err := server.Redis.SetValue(payload.ExecutionId, "execution failed"); err != nil {
			log.Printf("execution id: %s, failed to set value in redis cache - %v", payload.ExecutionId, err)
		}
		log.Printf("execution id: %s, failed to execute the code - %v", payload.ExecutionId, err)
		return
	}

	if result {
		if err := server.Redis.SetValue(payload.ExecutionId, "passed"); err != nil {
			log.Printf("execution id: %s, failed to set value in redis cache - %v", payload.ExecutionId, err)
		}
	} else {
		if err := server.Redis.SetValue(payload.ExecutionId, "failed"); err != nil {
			log.Printf("execution id: %s, failed to set value in redis cache - %v", payload.ExecutionId, err)
		}
	}
}
