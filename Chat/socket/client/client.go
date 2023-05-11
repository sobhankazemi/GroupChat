package client

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sobhankazemi/chat/models"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func Read(c *models.Client, q *amqp.Queue, ch *amqp.Channel) {
	defer func() {
		c.Hub.Disconnect <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, jsonmessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message := models.Message{}
		err = json.Unmarshal(jsonmessage, &message)
		if err != nil {
			log.Println("invalid json Request")
		}
		message.Room_id = c.Room_id
		message.Time = time.Now().Format("2006-01-02 15:04")
		go func(resp models.Message) {
			jsonResponse, _ := json.Marshal(resp)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = ch.PublishWithContext(ctx,
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "appliction/json",
					Body:        jsonResponse,
				})
			if err != nil {
				log.Println("faild to publish message")
			}
		}(message)
		c.Hub.BroadCast <- message
	}
}

func Write(c *models.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			jsonResponse, err := json.Marshal(message)
			if err != nil {
				log.Println("invalid json Response")
			}
			w.Write(jsonResponse)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				// w.Write(newline)
				temp := <-c.Send
				jsonResponse, err = json.Marshal(temp)
				if err != nil {
					log.Println("invalid response json")
				}
				w.Write(jsonResponse)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
