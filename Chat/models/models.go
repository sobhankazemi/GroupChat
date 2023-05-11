package models

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	Message  string `json:"message"`
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
	Time     string `json:"time"`
	Room_id  int `json:"room_id"`
}

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan Message
	Room_id int
}

type Hub struct {
	Clients    map[int]map[*Client]bool
	BroadCast  chan Message
	Connect    chan *Client
	Disconnect chan *Client
}
