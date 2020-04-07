package deck

import (
	"errors"
	"math/rand"
	"time"

	"github.com/peterzernia/set/ptr"
)

// Deck represents a deck in Set
type Deck struct {
	Cards []Card `json:"cards"`
}

// New creates a new deck
func New() *Deck {
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
	// deck.Shuffle()
	return &deck
}

// Shuffle uses Knuth shuffle algo to randomize the deck in O(n) time
// sourced from https://gist.github.com/quux00/8258425
func (d *Deck) Shuffle() {
	N := len(d.Cards)
	for i := 0; i < N; i++ {
		r := i + rand.Intn(N-i)
		d.Cards[r], d.Cards[i] = d.Cards[i], d.Cards[r]
	}
}

// CheckSet checks if 3 cards are a valid set. For each attribute,
// there are 3 options. If one of the attribute's options has a count of
// 2 that means the 3 cards are not a set, e.g. if there are two red cards
// the 3 cards are not a set.
func (d *Deck) CheckSet(cards []Card) (bool, error) {
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
