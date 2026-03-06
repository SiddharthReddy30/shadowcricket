package game

import (
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
)

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
