package game

import (
	"math/rand"
	"time"

	"github.com/peterzernia/set/ptr"
)

// Deck represents a deck in Set
type Deck struct {
	Cards []Card `json:"cards"`
}

// newDeck creates a new deck
func newDeck() *Deck {
	rand.Seed(time.Now().UnixNano())

	deck := Deck{}
	cards := []Card{}

	for _, color := range COLORS {
		for _, shape := range SHAPES {
			for _, number := range NUMBERS {
				for _, shading := range SHADINGS {
					card := Card{
						Color:   ptr.Int64(color),
						Shape:   ptr.Int64(shape),
						Number:  ptr.Int64(number),
						Shading: ptr.Int64(shading),
					}
					cards = append(cards, card)
				}
			}
		}
	}

	deck.Cards = cards
	deck.Shuffle()
	return &deck
}

// Shuffle uses Knuth shuffle algo to randomize the deck in O(n) time
// sourced from https://gist.github.com/quux00/8258425
func (d *Deck) Shuffle() {
	n := len(d.Cards)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		d.Cards[r], d.Cards[i] = d.Cards[i], d.Cards[r]
	}
}
