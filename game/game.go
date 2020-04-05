package game

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/peterzernia/set/deck"
	"github.com/peterzernia/set/ptr"
)

// Game represents a Set game
type Game struct {
	Deck     *deck.Deck    `json:"-"`
	GameOver *bool         `json:"game_over,omitempty"`
	InPlay   [][]deck.Card `json:"in_play,omitempty"`
	Players  []Player      `json:"players,omitempty"`
}

// New initializes a game
func New() *Game {
	game := Game{}
	game.Deck = deck.New()
	game.Deal()
	return &game
}

// Deal deals the initial 9 cards
func (g *Game) Deal() {
	inPlay := [][]deck.Card{[]deck.Card{}, []deck.Card{}, []deck.Card{}}
	inPlay[0] = g.Deck.Cards[0:4]
	inPlay[1] = g.Deck.Cards[4:8]
	inPlay[2] = g.Deck.Cards[8:12]

	g.InPlay = inPlay
	g.Deck.Cards = g.Deck.Cards[12:]
	return
}

// Play plays a play
func (g *Game) Play(move *Move, conn *websocket.Conn) (bool, error) {
	valid, err := g.Deck.CheckSet(move.Cards)
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

	// Replace the found set with new cards
	for _, v := range indices {
		if len(g.Deck.Cards) > 0 {
			g.InPlay[v[0]][v[1]] = g.Deck.Cards[0]
			g.Deck.Cards = g.Deck.Cards[1:]
		} else {
			g.InPlay[v[0]][v[1]] = deck.Card{}
		}
	}

	// Give the player a point
	g.UpdateScore(conn, 1)

	// If there are no cards left, check if there are any
	// remaining sets on the board, if not the game is over
	if len(g.Deck.Cards) == 0 {
		cards := g.InPlay[0]
		cards = append(cards, g.InPlay[1]...)
		cards = append(cards, g.InPlay[2]...)

		notEnd := false
		for i, x := range cards {
			for j, y := range cards {
				for k, z := range cards {
					if i != j && i != k && j == k &&
						x.Color != nil &&
						y.Color != nil &&
						z.Color != nil {
						valid, _ := g.Deck.CheckSet([]deck.Card{x, y, z})
						if !valid {
							notEnd = true
						}
					}
				}
			}
		}

		if !notEnd {
			return true, nil
		}
	}

	return false, nil
}

// AddCards adds another 3 cards onto the game board
// when there are no more sets left
func (g *Game) AddCards() {
	// TODO only add set when all players have requested
	for i := range g.InPlay {
		g.InPlay[i] = append(g.InPlay[i], g.Deck.Cards[0])
		g.Deck.Cards = g.Deck.Cards[1:]
	}
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
