package game

import "testing"

var testSecret = "test-secret-key-32-bytes-long!x!"

func TestCreateToken_ReturnsNonEmpty(t *testing.T) {
	token, err := CreateToken(1, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}

func TestCreateToken_UniqueTokens(t *testing.T) {
	token1, err := CreateToken(1, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	token2, err := CreateToken(1, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token1 == token2 {
		t.Error("expected different tokens for same player ID (unique nonce)")
	}
}

func TestDecryptToken_ReturnsCorrectPlayerID(t *testing.T) {
	token, err := CreateToken(42, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	playerID, err := DecryptToken(token, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if playerID != 42 {
		t.Errorf("expected player ID 42, got %d", playerID)
	}
}

func TestDecryptToken_WrongSecret(t *testing.T) {
	token, err := CreateToken(1, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = DecryptToken(token, "wrong-secret-key-32-bytes-long!!")
	if err == nil {
		t.Error("expected error when decrypting with wrong secret")
	}
}

func TestDecryptToken_TamperedToken(t *testing.T) {
	token, err := CreateToken(1, testSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tampered := token + "x"
	_, err = DecryptToken(tampered, testSecret)
	if err == nil {
		t.Error("expected error for tampered token")
	}
}

func TestDecryptToken_InvalidBase64(t *testing.T) {
	_, err := DecryptToken("!!!not-valid-base64!!!", testSecret)
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}
