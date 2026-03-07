package game

import (
	"fmt"
	"testing"

	"github.com/siddharthreddy/shadowcricket/internal/player"
)

var (
	virat = player.Player{
		ID: 1, Name: "Virat Kohli", Country: "India",
		JerseyNumber: 18, Role: "Middle-Order Batsman",
		IPLTeam: "Royal Challengers Bengaluru", IsWicketKeeper: false,
	}
	dhoni = player.Player{
		ID: 2, Name: "MS Dhoni", Country: "India",
		JerseyNumber: 7, Role: "Finisher",
		IPLTeam: "Chennai Super Kings", IsWicketKeeper: true,
	}
	abd = player.Player{
		ID: 3, Name: "AB de Villiers", Country: "South Africa",
		JerseyNumber: 17, Role: "Middle-Order Batsman",
		IPLTeam: "Royal Challengers Bengaluru", IsWicketKeeper: true,
	}
	bumrah = player.Player{
		ID: 4, Name: "Jasprit Bumrah", Country: "India",
		JerseyNumber: 93, Role: "Bowler",
		IPLTeam: "Mumbai Indians", IsWicketKeeper: false,
	}
	rashid = player.Player{
		ID: 5, Name: "Rashid Khan", Country: "Afghanistan",
		JerseyNumber: 19, Role: "Bowler",
		IPLTeam: "Gujarat Titans", IsWicketKeeper: false,
	}
	pant = player.Player{
		ID: 6, Name: "Rishabh Pant", Country: "India",
		JerseyNumber: 17, Role: "Finisher",
		IPLTeam: "Delhi Capitals", IsWicketKeeper: true,
	}
	rohit = player.Player{
		ID: 7, Name: "Rohit Sharma", Country: "India",
		JerseyNumber: 45, Role: "Opening Batsman",
		IPLTeam: "Mumbai Indians", IsWicketKeeper: false,
	}
	stokes = player.Player{
		ID: 8, Name: "Ben Stokes", Country: "England",
		JerseyNumber: 55, Role: "All-Rounder",
		IPLTeam: "Chennai Super Kings", IsWicketKeeper: false,
	}
	williamson = player.Player{
		ID: 9, Name: "Kane Williamson", Country: "New Zealand",
		JerseyNumber: 22, Role: "Middle-Order Batsman",
		IPLTeam: "", IsWicketKeeper: false,
	}
)

func assertField(t *testing.T, result GuessResult, field string, expected Color) {
	t.Helper()
	for _, f := range result.Feedback {
		if f.Field == field {
			if f.Color != expected {
				t.Errorf("field %s: expected %s, got %s", field, expected, f.Color)
			}
			return
		}
	}
	t.Errorf("field %s not found in feedback", field)
}

func assertFieldValue(t *testing.T, result GuessResult, field string, expected string) {
	t.Helper()
	for _, f := range result.Feedback {
		if f.Field == field {
			if f.Value != expected {
				t.Errorf("field %s: expected value %q, got %q", field, expected, f.Value)
			}
			return
		}
	}
	t.Errorf("field %s not found in feedback", field)
}

func printResult(label string, result GuessResult) {
	fmt.Printf("\n=== %s ===\n", label)
	fmt.Printf("Correct: %v\n", result.Correct)
	for _, f := range result.Feedback {
		fmt.Printf("  %-15s %-30s %s\n", f.Field, f.Value, f.Color)
	}
}
