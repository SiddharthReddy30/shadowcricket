package game

type Color string

const (
	Green  Color = "green"
	Yellow Color = "yellow"
	White  Color = "white"
)

type FieldFeedback struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Color Color  `json:"color"`
}

type GuessResult struct {
	Correct  bool            `json:"correct"`
	Feedback []FieldFeedback `json:"feedback"`
}

type tokenPayload struct {
	PlayerID int `json:"pid"`
}

var countryToContinent = map[string]string{
	"India":        "Asia",
	"Afghanistan":  "Asia",
	"Pakistan":     "Asia",
	"Sri Lanka":    "Asia",
	"Bangladesh":   "Asia",
	"Australia":    "Oceania",
	"New Zealand":  "Oceania",
	"England":      "Europe",
	"South Africa": "Africa",
	"West Indies":  "Americas",
	"Zimbabwe":     "Africa",
	"USA":          "Americas",
}

var roleRank = map[string]int{
	"Opening Batsman":      1,
	"Middle-Order Batsman": 2,
	"Finisher":             3,
	"All-Rounder":          4,
	"Bowler":               5,
}
