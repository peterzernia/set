package deck

// Constants for colors
const (
	RED int64 = iota
	BLUE
	GREEN
)

// Constants for shapes
const (
	DIAMOND int64 = iota
	OVAL
	SQUIGGLE
)

// Constants for numbers
const (
	ONE int64 = iota
	TWO
	THREE
)

// Constants for shading
const (
	OUTLINED int64 = iota
	STRIPED
	SOLID
)

// Global variables representing default colors, shapes and numbers
// for the cards in a Set deck
var (
	COLORS   = []int64{RED, BLUE, GREEN}
	SHAPES   = []int64{DIAMOND, OVAL, SQUIGGLE}
	NUMBERS  = []int64{ONE, TWO, THREE}
	SHADINGS = []int64{OUTLINED, STRIPED, SOLID}
)

// Card represents a card in Set
type Card struct {
	Color   *int64 `json:"color"`
	Shape   *int64 `json:"shape"`
	Number  *int64 `json:"number"`
	Shading *int64 `json:"shading"`
}
