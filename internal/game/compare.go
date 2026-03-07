package game

import (
	"strings"

	"github.com/siddharthreddy/shadowcricket/internal/player"
)

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
