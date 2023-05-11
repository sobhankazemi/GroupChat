package main

import (
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sobhankazemi/chat/handlers"
	"github.com/sobhankazemi/chat/socket/hub"
)

var handler *handlers.Repository

func main() {

	myhub := hub.NewHub()
	go hub.StartHub(myhub)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("can not connect to rabbitmq service")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("can not open up a rabbit channel")
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"chat-history",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("can not initialize a rabbit queue")
	}
	handler = handlers.NewHandler(myhub, &q, ch)
	server := http.Server{
		Addr:    ":8080",
		Handler: GetRoutes(),
	}
	fmt.Println("application running on port 8080 ...")
	server.ListenAndServe()
}
