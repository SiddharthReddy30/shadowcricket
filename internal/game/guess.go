package game

import (
	"fmt"

	"github.com/siddharthreddy/shadowcricket/internal/player"
)

func EvaluateGuess(target, guess player.Player) GuessResult {
	feedback := []FieldFeedback{
		{Field: "name", Value: guess.Name, Color: compareName(target.Name, guess.Name)},
		{Field: "country", Value: guess.Country, Color: compareCountry(target.Country, guess.Country)},
		{Field: "jersey_number", Value: fmt.Sprintf("%d", guess.JerseyNumber), Color: compareJerseyNumber(target.JerseyNumber, guess.JerseyNumber)},
		compareRole(target, guess),
		{Field: "ipl_team", Value: guess.IPLTeam, Color: compareIPLTeam(target.IPLTeam, guess.IPLTeam)},
	}
	return GuessResult{
		Correct:  target.ID == guess.ID,
		Feedback: feedback,
	}
}
