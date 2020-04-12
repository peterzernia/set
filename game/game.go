package game

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/set/deck"
	"github.com/peterzernia/set/ptr"
)

// Game represents a Set game
type Game struct {
	Deck       *deck.Deck    `json:"-"`
	GameOver   *bool         `json:"game_over,omitempty"`
	InPlay     [][]deck.Card `json:"in_play,omitempty"`
	LastPlayer *string       `json:"last_player,omitempty"`
	LastSet    []deck.Card   `json:"last_set,omitempty"`
	Players    []Player      `json:"players,omitempty"`
	Remaining  *int64        `json:"remaining,omitempty"`
}

// New initializes a game
func New() *Game {
	game := Game{}
	game.Players = []Player{}
	game.Deck = deck.New()
	game.Deal()
	return &game
}

// Deal deals the initial 12 cards
func (g *Game) Deal() {
	inPlay := [][]deck.Card{[]deck.Card{}, []deck.Card{}, []deck.Card{}}

	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			inPlay[j] = append(inPlay[j], *g.Deck.Cards[j].Copy())
		}
		g.Deck.Cards = g.Deck.Cards[3:]
	}

	g.InPlay = inPlay
	g.Remaining = ptr.Int64(int64(len(g.Deck.Cards)))
	return
}

// Play plays a play
func (g *Game) Play(move *Move, conn *websocket.Conn) (bool, error) {
	valid, err := g.CheckSet(move.Cards)
	if !valid || err != nil {
		g.UpdateScore(conn, -1)
		return false, errors.New("Invalid set")
	}

	// Find the location of the 3 cards in the cards in play matrix
	indices := [][]int{}
	for _, v := range move.Cards {
		i := findIndex(g.InPlay, v)
		indices = append(indices, i)
	}

	// -1 index signifies that the card is not in play
	for _, v := range indices {
		if v[0] == -1 || v[1] == -1 {
			return false, errors.New("Invalid cards")
		}
	}

	// Sum the cards in play, ignoring the placeholder ones
	inPlay := 0
	for _, row := range g.InPlay {
		for _, card := range row {
			if card.Color != nil {
				inPlay++
			}
		}
	}

	// Replace the found set with new cards if there aren't
	// 12 (+3 about to be removed) cards already in play
	for _, v := range indices {
		if len(g.Deck.Cards) > 0 && inPlay < 15 {
			g.InPlay[v[0]][v[1]] = *g.Deck.Cards[0].Copy()
			g.Deck.Cards = g.Deck.Cards[1:]
		} else {
			g.InPlay[v[0]][v[1]] = deck.Card{} // Placeholder card
		}
	}

	// Give the player a point
	g.UpdateScore(conn, 1)
	g.Remaining = ptr.Int64(int64(len(g.Deck.Cards)))

	var lastPlayer *string
	for _, player := range g.Players {
		if conn == player.Conn {
			lastPlayer = player.Name
		}
	}
	g.LastPlayer = lastPlayer
	g.LastSet = move.Cards

	// If there are no cards left, check if there are any
	// remaining sets on the board, if not the game is over
	if len(g.Deck.Cards) == 0 {
		end := g.CheckRemainingSets()

		return end, nil
	}

	return false, nil
}

// CheckSet checks if 3 cards are a valid set. For each attribute,
// there are 3 options. If one of the attribute's options has a count of
// 2 that means the 3 cards are not a set, e.g. if there are two red cards
// the 3 cards are not a set.
func (g *Game) CheckSet(cards []deck.Card) (bool, error) {
	if len(cards) != 3 {
		return false, errors.New("Sets must contain 3 cards")
	}

	colors := make(map[int64]int64)
	shapes := make(map[int64]int64)
	numbers := make(map[int64]int64)
	shadings := make(map[int64]int64)

	for _, card := range cards {
		colors[*card.Color]++
		shapes[*card.Shape]++
		numbers[*card.Number]++
		shadings[*card.Shading]++
	}

	for _, v := range colors {
		if v == 2 {
			return false, nil
		}
	}

	for _, v := range shapes {
		if v == 2 {
			return false, nil
		}
	}

	for _, v := range numbers {
		if v == 2 {
			return false, nil
		}
	}

	for _, v := range shadings {
		if v == 2 {
			return false, nil
		}
	}

	return true, nil
}

// CheckRemainingSets checks if there is still
// at least on set in the InPlay cards
func (g *Game) CheckRemainingSets() bool {
	cards := []deck.Card{}
	for _, row := range g.InPlay {
		for _, card := range row {
			if card.Color != nil {
				cards = append(cards, *card.Copy())
			}
		}
	}

	end := true
	for i, x := range cards {
		for j, y := range cards {
			for k, z := range cards {
				if i != j && i != k && j != k {
					valid, _ := g.CheckSet([]deck.Card{x, y, z})
					if valid {
						end = false
					}
				}
			}
		}
	}

	return end
}

// AddCards adds another 3 cards onto the game board
func (g *Game) AddCards() {
	cards := g.Deck.Cards[0:3]
	g.Deck.Cards = g.Deck.Cards[3:]

	// Fill in the spaces
	for i, row := range g.InPlay {
		for j, v := range row {
			if len(cards) > 0 {
				if v.Color == nil {
					g.InPlay[i][j] = cards[0]
					cards = cards[1:]
				}
			}
		}
	}

	// Start a new column if there are no more spaces
	if len(cards) > 0 {
		g.InPlay[0] = append(g.InPlay[0], *cards[0].Copy())
		g.InPlay[1] = append(g.InPlay[1], *cards[1].Copy())
		g.InPlay[2] = append(g.InPlay[2], *cards[2].Copy())
	}

	g.Remaining = ptr.Int64(int64(len(g.Deck.Cards)))
}

// UpdateScore updates a players score. The player is found by their websocket
// connection
func (g *Game) UpdateScore(conn *websocket.Conn, value int64) {
	var index int
	for i, player := range g.Players {
		if player.Conn == conn {
			index = i
		}
	}

	if *g.Players[index].Score == 0 && value < 0 {
		return
	}

	g.Players[index].Score = ptr.Int64(*g.Players[index].Score + value)
}
