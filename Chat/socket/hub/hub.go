package hub

import (
	"log"

	"github.com/sobhankazemi/chat/models"
)

func NewHub() *models.Hub {
	return &models.Hub{
		BroadCast:  make(chan models.Message),
		Connect:    make(chan *models.Client),
		Disconnect: make(chan *models.Client),
		Clients:    make(map[int]map[*models.Client]bool),
	}
}

func StartHub(h *models.Hub) {
	for {
		select {
		case client := <-h.Connect:
			if h.Clients[client.Room_id] == nil {
				h.Clients[client.Room_id] = make(map[*models.Client]bool)
			}
			h.Clients[client.Room_id][client] = true
			log.Printf("a client connected to room number %d\n", client.Room_id)
		case client := <-h.Disconnect:
			if _, ok := h.Clients[client.Room_id][client]; ok {
				log.Printf("a client disconnected from room number %d\n", client.Room_id)
				delete(h.Clients[client.Room_id], client)
				close(client.Send)
			}
		case message := <-h.BroadCast:
			for client := range h.Clients[message.Room_id] {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients[message.Room_id], client)
				}
			}
		}
	}
}
