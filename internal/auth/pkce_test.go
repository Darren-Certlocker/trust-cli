package auth

import "testing"

func TestChallengeMatchesRFC7636(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	challenge := Challenge(verifier)
	if challenge != "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM" {
		t.Fatalf("challenge %s", challenge)
	}
}

func TestNewVerifierLength(t *testing.T) {
	verifier, err := NewVerifier()
	if err != nil {
		t.Fatal(err)
	}
	if len(verifier) < 43 || len(verifier) > 128 {
		t.Fatalf("verifier length %d", len(verifier))
	}
	if Challenge(verifier) == "" {
		t.Fatal("empty challenge")
	}
}
