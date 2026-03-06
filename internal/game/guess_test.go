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

func TestEvaluateGuess_ExactMatch(t *testing.T) {
	result := EvaluateGuess(virat, virat)
	if !result.Correct {
		t.Error("expected correct=true for same player")
	}
	for _, f := range result.Feedback {
		if f.Color != Green {
			t.Errorf("expected green for all fields, got %s for %s", f.Color, f.Field)
		}
	}
	printResult("Target: Virat, Guess: Virat", result)
}

func TestEvaluateGuess_SameCountryDifferentPlayer(t *testing.T) {
	result := EvaluateGuess(virat, bumrah)
	if result.Correct {
		t.Error("expected correct=false")
	}
	printResult("Target: Virat, Guess: Bumrah", result)

	// Country should be green (both India)
	assertField(t, result, "country", Green)
	// Name should be white
	assertField(t, result, "name", White)
}

func TestEvaluateGuess_SameContinent(t *testing.T) {
	// India and Afghanistan are both Asia
	result := EvaluateGuess(virat, rashid)
	printResult("Target: Virat, Guess: Rashid", result)

	assertField(t, result, "country", Yellow)
}

func TestEvaluateGuess_DifferentContinent(t *testing.T) {
	// India and South Africa
	result := EvaluateGuess(virat, abd)
	printResult("Target: Virat, Guess: ABD", result)

	assertField(t, result, "country", White)
}

func TestEvaluateGuess_JerseyNumberExact(t *testing.T) {
	// ABD (17) vs Pant (17)
	result := EvaluateGuess(abd, pant)
	printResult("Target: ABD, Guess: Pant", result)

	assertField(t, result, "jersey_number", Green)
}

func TestEvaluateGuess_JerseyNumberClose(t *testing.T) {
	// Virat (18) vs Rashid (19) — diff 1
	result := EvaluateGuess(virat, rashid)
	printResult("Target: Virat, Guess: Rashid", result)

	assertField(t, result, "jersey_number", Yellow)
}

func TestEvaluateGuess_JerseyNumberFar(t *testing.T) {
	// Virat (18) vs Bumrah (93) — diff 75
	result := EvaluateGuess(virat, bumrah)
	printResult("Target: Virat, Guess: Bumrah", result)

	assertField(t, result, "jersey_number", White)
}

func TestEvaluateGuess_BothWicketKeepers(t *testing.T) {
	// Dhoni (Finisher, WK) vs Pant (Finisher, WK)
	result := EvaluateGuess(dhoni, pant)
	printResult("Target: Dhoni, Guess: Pant", result)

	assertField(t, result, "role", Green)
	assertFieldValue(t, result, "role", "Wicket-Keeper")
}

func TestEvaluateGuess_OneWicketKeeper(t *testing.T) {
	// Dhoni (Finisher, WK) vs Virat (Middle-Order, not WK)
	result := EvaluateGuess(dhoni, virat)
	printResult("Target: Dhoni, Guess: Virat", result)

	// Should compare positions: Finisher(3) vs Middle-Order(2) = diff 1 = yellow
	assertField(t, result, "role", Yellow)
	assertFieldValue(t, result, "role", "Middle-Order Batsman")
}

func TestEvaluateGuess_RoleFarApart(t *testing.T) {
	// Rohit (Opening, 1) vs Bumrah (Bowler, 5) — diff 4
	result := EvaluateGuess(rohit, bumrah)
	printResult("Target: Rohit, Guess: Bumrah", result)

	assertField(t, result, "role", White)
}

func TestEvaluateGuess_RoleOneApart(t *testing.T) {
	// Stokes (All-Rounder, 4) vs Bumrah (Bowler, 5) — diff 1
	result := EvaluateGuess(stokes, bumrah)
	printResult("Target: Stokes, Guess: Bumrah", result)

	assertField(t, result, "role", Yellow)
}

func TestEvaluateGuess_SameIPLTeam(t *testing.T) {
	// Virat and ABD both RCB
	result := EvaluateGuess(virat, abd)
	printResult("Target: Virat, Guess: ABD", result)

	assertField(t, result, "ipl_team", Green)
}

func TestEvaluateGuess_DifferentIPLTeam(t *testing.T) {
	// Virat (RCB) vs Dhoni (CSK) — both in IPL
	result := EvaluateGuess(virat, dhoni)
	printResult("Target: Virat, Guess: Dhoni", result)

	assertField(t, result, "ipl_team", Yellow)
}

func TestEvaluateGuess_NoIPLTeam(t *testing.T) {
	// Virat (RCB) vs Williamson (no IPL team)
	result := EvaluateGuess(virat, williamson)
	printResult("Target: Virat, Guess: Williamson", result)

	assertField(t, result, "ipl_team", White)
}

// --- Helpers ---

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
