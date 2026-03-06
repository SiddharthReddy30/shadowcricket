package game

import (
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"github.com/siddharthreddy/shadowcricket/internal/player"
)

// --- Token ---

func CreateToken(playerID int, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := cryptoRand.Read(nonce); err != nil {
		return "", err
	}
	payload, err := json.Marshal(tokenPayload{PlayerID: playerID})
	if err != nil {
		return "", err
	}
	encrypted := gcm.Seal(nonce, nonce, payload, nil)
	return base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func DecryptToken(token, secret string) (int, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, errors.New("invalid token")
	}
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return 0, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return 0, errors.New("invalid token")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	payload, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, errors.New("invalid or tampered token")
	}
	var p tokenPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return 0, err
	}
	return p.PlayerID, nil
}

// --- Guess Evaluation ---

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
