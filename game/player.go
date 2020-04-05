package game

import "github.com/gorilla/websocket"

// Player represents a player
type Player struct {
	ID      *int64          `json:"id,omitempty"`
	Conn    *websocket.Conn `json:"-"`
	Name    *string         `json:"name,omitempty"`
	Score   *int64          `json:"score,omitempty"`
	Request *bool           `json:"request,omitempty"`
}
