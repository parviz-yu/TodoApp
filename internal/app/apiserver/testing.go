package apiserver

import (
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func TestSession(t *testing.T) (*sessions.CookieStore, *securecookie.SecureCookie) {
	t.Helper()

	secretKey := []byte("secret")
	cookieStore := sessions.NewCookieStore(secretKey)
	secureCookie := securecookie.New(secretKey, nil)

	return cookieStore, secureCookie
}
