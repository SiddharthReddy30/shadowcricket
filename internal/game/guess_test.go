package game

import "testing"

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

func TestEvaluateGuess_DifferentPlayer(t *testing.T) {
	result := EvaluateGuess(virat, bumrah)
	if result.Correct {
		t.Error("expected correct=false")
	}
	printResult("Target: Virat, Guess: Bumrah", result)

	assertField(t, result, "country", Green)
	assertField(t, result, "name", White)
}
