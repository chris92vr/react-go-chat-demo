package websocket

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type     int    `json:"type"`
	Time     string `json:"time"`
	Body     string `json:"body"`
	ClientID string `json:"client_id"`
}

func generateRandomString() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		now := time.Now().Format("15:04:05 02-01-2006 ")
		c.ID = generateRandomString()

		client := &Client{ID: c.ID, Conn: c.Conn, Pool: c.Pool}
		message := Message{Type: messageType, Time: now, Body: string(p), ClientID: c.ID}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
		fmt.Printf("Message Received: %+v\n", client.ID)

	}
}
