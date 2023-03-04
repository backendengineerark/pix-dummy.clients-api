package handlers

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/backendengineerark/clients-api/pkg/customlogs"
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

func (ach *AccountCreatedNotifyHandler) Handle(ctx context.Context, event events.EventInterface, wg *sync.WaitGroup) {
	logger := customlogs.GetContextLogger(ctx)

	logger.Printf("Try to notify by rabbimq")
	defer wg.Done()
	jsonOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		logger.Printf("fail to encode payload to json %s", event.GetPayload())
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
		logger.Printf("fail to send message to rabbimq beacause %s", err)
	}

	logger.Printf("Success to notify rabbitmq")
}
