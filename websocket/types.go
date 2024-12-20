package websocket

import "github.com/gorilla/websocket"

type Client struct {
	Conn     *websocket.Conn
	DeviceID string
	Room     string
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
