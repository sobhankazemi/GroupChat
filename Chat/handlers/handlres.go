package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sobhankazemi/chat/models"
	"github.com/sobhankazemi/chat/socket/client"
)

type Repository struct {
	Queue    *amqp.Queue
	Channel  *amqp.Channel
	Hub      *models.Hub
	upgrader websocket.Upgrader
}

func NewHandler(hub *models.Hub, q *amqp.Queue, ch *amqp.Channel) *Repository {
	return &Repository{
		Queue:   q,
		Channel: ch,
		Hub:     hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (repo *Repository) ServeChatRooms(w http.ResponseWriter, r *http.Request) {
	//authentication
	//decode jwt
	conn, err := repo.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// token := r.URL.Query().Get("token")
	// claims := jwt.MapClaims{}
	// _, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("django-insecure-8e-v!7wefcxi8yd^p+8!d%mh2lf#796mo71b5u9@l9(eh@)$i_"), nil
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	conn.WriteMessage(websocket.CloseMessage, []byte("not authenticated"))
	// 	return
	// }
	// username := claims["username"]
	// user_id := claims["user_id"]
	room_id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	user := &models.Client{Hub: repo.Hub, Conn: conn, Send: make(chan models.Message), Room_id: room_id}
	user.Hub.Connect <- user

	go client.Write(user)
	go client.Read(user, repo.Queue , repo.Channel)
}
