package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Pool *Pool
	ID   string
}

type Message struct {
	Type       int    `json:"type"`
	Time       string `json:"time"`
	Body       string `json:"body"`
	Clientname string `json:"client_name"`
}

func generateRandomString() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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

		client := &Client{Conn: c.Conn, Pool: c.Pool, ID: c.ID}
		message := Message{Type: messageType, Time: now, Body: string(p), Clientname: c.ID}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
		fmt.Printf("Message Received: %+v\n", client)

	}
}
