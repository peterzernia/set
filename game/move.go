package game

// Move represents a play made by a player
type Move struct {
	PlayerID *int64 `json:"player_id"`
	Cards    []Card `json:"cards"`
}
