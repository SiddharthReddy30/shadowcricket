package game

import (
	"fmt"
	"strings"

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

func compareName(target, guess string) Color {
	if strings.EqualFold(target, guess) {
		return Green
	}
	return White
}

func compareCountry(target, guess string) Color {
	if strings.EqualFold(target, guess) {
		return Green
	}
	if countryToContinent[target] == countryToContinent[guess] {
		return Yellow
	}
	return White
}

func compareJerseyNumber(target, guess int) Color {
	if target == guess {
		return Green
	}
	diff := target - guess
	if diff < 0 {
		diff = -diff
	}
	if diff <= 5 {
		return Yellow
	}
	return White
}

func compareRole(target, guess player.Player) FieldFeedback {
	if target.IsWicketKeeper && guess.IsWicketKeeper {
		return FieldFeedback{Field: "role", Value: "Wicket-Keeper", Color: Green}
	}
	diff := roleRank[target.Role] - roleRank[guess.Role]
	if diff < 0 {
		diff = -diff
	}
	color := White
	if diff == 0 {
		color = Green
	} else if diff == 1 {
		color = Yellow
	}
	return FieldFeedback{Field: "role", Value: guess.Role, Color: color}
}

func compareIPLTeam(target, guess string) Color {
	if strings.EqualFold(target, guess) {
		return Green
	}
	if target != "" && guess != "" {
		return Yellow
	}
	return White
}
