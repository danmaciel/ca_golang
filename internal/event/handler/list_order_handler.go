package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/danmaciel/ca_golang/pkg/events"
	"github.com/streadway/amqp"
)

type ListOrderHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewListOrderHandler(rabbitMQChannel *amqp.Channel) *ListOrderHandler {
	return &ListOrderHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *ListOrderHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)
}
