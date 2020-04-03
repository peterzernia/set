package game

// Player represents a player
type Player struct {
	ID      *int64  `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
	Score   *int64  `json:"score,omitempty"`
	Request *bool   `json:"request,omitempty"`
}
