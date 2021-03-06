package game

// findIndex is a helper function to find the index of a card
func findIndex(cardss [][]Card, card Card) []int {
	for i, cards := range cardss {
		for j, v := range cards {
			if v.Color != nil && *v.Color == *card.Color &&
				*v.Shape == *card.Shape &&
				*v.Number == *card.Number &&
				*v.Shading == *card.Shading {
				return []int{i, j}
			}
		}
	}
	return []int{-1, -1}
}
