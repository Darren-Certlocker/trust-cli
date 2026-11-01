package auth

import "testing"

func TestLoginURL(t *testing.T) {
	got, err := LoginURL("https://app.certlocker.io/", "http://127.0.0.1:9/callback", "state", "challenge")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://app.certlocker.io/login?code_challenge=challenge&redirect_uri=http%3A%2F%2F127.0.0.1%3A9%2Fcallback&state=state"
	if got != want {
		t.Fatalf("got %s", got)
	}
}
