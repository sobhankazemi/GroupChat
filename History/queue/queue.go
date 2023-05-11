package queue

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sobhankazemi/GroupChat/History/dbrepo"
	"github.com/sobhankazemi/GroupChat/History/models"
)

type MQueue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func InitQueue(connectionString, channelName string) (*MQueue, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		channelName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	queue := &MQueue{
		conn: conn,
		ch:   ch,
		q:    q,
	}
	return queue, nil
}

func (queue *MQueue) Listen(postgreRepo dbrepo.Repository) error {
	defer queue.conn.Close()
	defer queue.ch.Close()
	msgs, err := queue.ch.Consume(
		queue.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for d := range msgs {
		message := models.Message{}
		err := json.Unmarshal(d.Body, &message)
		if err != nil {
			log.Println("invalid json message")
			continue
		}
		postgreRepo.SaveChatHistory(message)
	}
	return nil
}
