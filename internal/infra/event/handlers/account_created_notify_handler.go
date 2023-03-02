package handlers

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/backendengineerark/clients-api/pkg/events"
	"github.com/streadway/amqp"
)

type AccountCreatedNotifyHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewAccountCreatedNotifyHandler(rabbitMQChannel *amqp.Channel) *AccountCreatedNotifyHandler {
	return &AccountCreatedNotifyHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (ach *AccountCreatedNotifyHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Try to notify")
	jsonOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		fmt.Println("fail to encode")
	}

	err = ach.RabbitMQChannel.Publish(
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonOutput,
		},
	)
	if err != nil {
		fmt.Println("fail to send")
	}

	fmt.Println("Success to notify rabbitmq")
}
