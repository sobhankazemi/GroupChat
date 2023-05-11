package main

import (
	"fmt"
	"log"
	"net/http"

	postgresrepo "github.com/sobhankazemi/GroupChat/History/dbrepo/postgresRepo"
	"github.com/sobhankazemi/GroupChat/History/driver"
	handlers "github.com/sobhankazemi/GroupChat/History/handlres"
	"github.com/sobhankazemi/GroupChat/History/queue"
)

var handler *handlers.Repository

func main() {
	db, err := driver.NewConnection("host=db port=5432 dbname=temp_chat user=postgres password=1234")
	if err != nil {
		log.Fatal("can not connect to database")
	}
	posrgreRepo := postgresrepo.NewPostgreRepo(db)
	handler = handlers.NewHandler(posrgreRepo)
	rabbitQueue, err := queue.InitQueue("amqp://guest:guest@rabbit:5672/", "chat-history")
	if err != nil {
		log.Fatal("faild to set up rabbitMQ")
	}
	go rabbitQueue.Listen(posrgreRepo)

	server := http.Server{
		Addr:    ":8050",
		Handler: GetRoutes(),
	}
	fmt.Println("listens on port 8050")
	server.ListenAndServe()
}
