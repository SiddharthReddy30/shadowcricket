package game

import "testing"

func TestCompareCountry_SameCountry(t *testing.T) {
	result := EvaluateGuess(virat, bumrah)
	printResult("Target: Virat, Guess: Bumrah", result)

	assertField(t, result, "country", Green)
}

func TestCompareCountry_SameContinent(t *testing.T) {
	// India and Afghanistan are both Asia
	result := EvaluateGuess(virat, rashid)
	printResult("Target: Virat, Guess: Rashid", result)

	assertField(t, result, "country", Yellow)
}

func TestCompareCountry_DifferentContinent(t *testing.T) {
	// India and South Africa
	result := EvaluateGuess(virat, abd)
	printResult("Target: Virat, Guess: ABD", result)

	assertField(t, result, "country", White)
}

func TestCompareJerseyNumber_Exact(t *testing.T) {
	// ABD (17) vs Pant (17)
	result := EvaluateGuess(abd, pant)
	printResult("Target: ABD, Guess: Pant", result)

	assertField(t, result, "jersey_number", Green)
}

func TestCompareJerseyNumber_Close(t *testing.T) {
	// Virat (18) vs Rashid (19) — diff 1
	result := EvaluateGuess(virat, rashid)
	printResult("Target: Virat, Guess: Rashid", result)

	assertField(t, result, "jersey_number", Yellow)
}

func TestCompareJerseyNumber_Far(t *testing.T) {
	// Virat (18) vs Bumrah (93) — diff 75
	result := EvaluateGuess(virat, bumrah)
	printResult("Target: Virat, Guess: Bumrah", result)

	assertField(t, result, "jersey_number", White)
}

func TestCompareRole_BothWicketKeepers(t *testing.T) {
	// Dhoni (Finisher, WK) vs Pant (Finisher, WK)
	result := EvaluateGuess(dhoni, pant)
	printResult("Target: Dhoni, Guess: Pant", result)

	assertField(t, result, "role", Green)
	assertFieldValue(t, result, "role", "Wicket-Keeper")
}

func TestCompareRole_OneWicketKeeper(t *testing.T) {
	// Dhoni (Finisher, WK) vs Virat (Middle-Order, not WK)
	result := EvaluateGuess(dhoni, virat)
	printResult("Target: Dhoni, Guess: Virat", result)

	assertField(t, result, "role", Yellow)
	assertFieldValue(t, result, "role", "Middle-Order Batsman")
}

func TestCompareRole_FarApart(t *testing.T) {
	// Rohit (Opening, 1) vs Bumrah (Bowler, 5) — diff 4
	result := EvaluateGuess(rohit, bumrah)
	printResult("Target: Rohit, Guess: Bumrah", result)

	assertField(t, result, "role", White)
}

func TestCompareRole_OneApart(t *testing.T) {
	// Stokes (All-Rounder, 4) vs Bumrah (Bowler, 5) — diff 1
	result := EvaluateGuess(stokes, bumrah)
	printResult("Target: Stokes, Guess: Bumrah", result)

	assertField(t, result, "role", Yellow)
}

func TestCompareIPLTeam_SameTeam(t *testing.T) {
	// Virat and ABD both RCB
	result := EvaluateGuess(virat, abd)
	printResult("Target: Virat, Guess: ABD", result)

	assertField(t, result, "ipl_team", Green)
}

func TestCompareIPLTeam_DifferentTeam(t *testing.T) {
	// Virat (RCB) vs Dhoni (CSK) — both in IPL
	result := EvaluateGuess(virat, dhoni)
	printResult("Target: Virat, Guess: Dhoni", result)

	assertField(t, result, "ipl_team", Yellow)
}

func TestCompareIPLTeam_NoTeam(t *testing.T) {
	// Virat (RCB) vs Williamson (no IPL team)
	result := EvaluateGuess(virat, williamson)
	printResult("Target: Virat, Guess: Williamson", result)

	assertField(t, result, "ipl_team", White)
}
