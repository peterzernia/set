package game

import (
	"errors"
	"fmt"

	"github.com/peterzernia/set/deck"
	"github.com/peterzernia/set/ptr"
)

// Game represents a Set game
type Game struct {
	Deck    *deck.Deck    `json:"-"`
	InPlay  [][]deck.Card `json:"in_play,omitempty"`
	Players []Player      `json:"players,omitempty"`
	Winner  *Player       `json:"winner,omitempty"`
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
func (g *Game) Play(move *Move) error {
	valid, err := g.Deck.CheckSet(move.Cards)
	if !valid || err != nil {
		return errors.New("Invalid set")
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
			return errors.New("Invalid cards")
		}
	}

	// Replace the found set with new cards
	for _, v := range indices {
		if len(g.Deck.Cards) > 0 {
			fmt.Println(len(g.Deck.Cards))
			g.Deck.Cards = g.Deck.Cards[1:]
			g.InPlay[v[0]][v[1]] = g.Deck.Cards[0]
		} else {
			g.InPlay[v[0]][v[1]] = deck.Card{}
		}
	}

	// Give the player a point
	g.Players[*move.PlayerID-1].Score = ptr.Int64(*g.Players[*move.PlayerID-1].Score + 1)

	return nil
}

// Refresh deals out cards into the places
// the last set was
func (g *Game) Refresh() {
	// TODO only refresh when there are fewer than 12 cards in play
	for i, row := range g.InPlay {
		for j, card := range row {
			if card.Color == nil {
				g.InPlay[i][j] = g.Deck.Cards[0]
				g.Deck.Cards = g.Deck.Cards[1:]
			}
		}
	}
	return
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
