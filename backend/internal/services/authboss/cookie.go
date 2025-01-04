package authboss

import (
	"net/http"
	"strings"
	"time"

	"github.com/volatiletech/authboss/v3"
)

// Cookie storer implements authboss.ClientStateReadWriter for cookies
type CookieStorer struct {
	secure   bool
	sameSite http.SameSite
	domain   string
}

// CookieConfig holds the configuration for the CookieStorer
type CookieConfig struct {
	Secure   bool
	SameSite http.SameSite
	Domain   string
}

// NewCookieStorer creates a new CookieStorer with the given configuration
func NewCookieStorer(cfg CookieConfig) *CookieStorer {
	return &CookieStorer{
		secure:   cfg.Secure,
		sameSite: cfg.SameSite,
		domain:   cfg.Domain,
	}
}

// CookieState implements authboss.ClientState
type CookieState struct {
	cookies map[string]string
}

// Get retrieves a cookie value
func (s *CookieState) Get(key string) (string, bool) {
	if s.cookies == nil {
		return "", false
	}
	val, ok := s.cookies[key]
	return val, ok
}

// Set stores a cookie value
func (s *CookieState) Set(key, val string) {
	if s.cookies == nil {
		s.cookies = make(map[string]string)
	}
	s.cookies[key] = val
}

// Delete removes a cookie
func (s *CookieState) Delete(key string) {
	if s.cookies != nil {
		delete(s.cookies, key)
	}
}

// ReadState loads the cookie state from the request
func (s *CookieStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	state := &CookieState{
		cookies: make(map[string]string),
	}

	for _, cookie := range r.Cookies() {
		state.Set(cookie.Name, cookie.Value)
	}

	return state, nil
}

// WriteState writes the cookie state to the response
func (s *CookieStorer) WriteState(w http.ResponseWriter, r *http.Request, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	for _, ev := range events {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			http.SetCookie(w, &http.Cookie{
				Name:     ev.Key,
				Value:    ev.Value,
				Path:     "/",
				Domain:   s.domain,
				Secure:   s.secure,
				HttpOnly: true,
				SameSite: s.sameSite,
				// Remember me cookie lasts 30 days
				MaxAge: int(30 * 24 * time.Hour.Seconds()), // 30 days
			})

		case authboss.ClientStateEventDel:
			http.SetCookie(w, &http.Cookie{
				Name:     ev.Key,
				Value:    "",
				Path:     "/",
				Domain:   s.domain,
				Secure:   s.secure,
				HttpOnly: true,
				SameSite: s.sameSite,
				MaxAge:   -1,
			})

		case authboss.ClientStateEventDelAll:
			whitelist := make(map[string]bool)
			for _, k := range strings.Split(ev.Key, ",") {
				whitelist[strings.TrimSpace(k)] = true
			}

			for _, cookie := range r.Cookies() {
				if !whitelist[cookie.Name] {
					http.SetCookie(w, &http.Cookie{
						Name:     cookie.Name,
						Value:    "",
						Path:     "/",
						Domain:   s.domain,
						Secure:   s.secure,
						HttpOnly: true,
						SameSite: s.sameSite,
						MaxAge:   -1,
					})
				}
			}
		}
	}

	return nil
}
