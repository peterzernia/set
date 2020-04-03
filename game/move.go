package game

import "github.com/peterzernia/set/deck"

// Move represents a play made by a player
type Move struct {
	PlayerID *int64      `json:"player_id"`
	Cards    []deck.Card `json:"cards"`
}
