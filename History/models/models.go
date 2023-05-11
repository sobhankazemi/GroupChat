package models

type Message struct {
	Message  string `json:"message"`
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
	Time     string `json:"time"`
	Room_id  int `json:"room_id"`
}
