package main

import (
	"context"
	"gopkg.in/aabizri/goil.v0"
	"net/http"
	"strings"
)

const (
	AuthErrorAuthorizationFailure                 string = "A0"
	AuthErrorAuthorizationHeaderAbsent                   = "A10"
	AuthErrorAuthorizationHeaderInvalid                  = "A11"
	AuthErrorAuthorizationHeaderElementsLenNotTwo        = "A110"
	AuthErrorAuthorizationHeaderBearerTagAbsent          = "A111"
	AuthErrorInvalidToken                                = "A20"
)

// authenticationMiddleware wraps the normal handler around in order to satisfy http.Handler type
type authenticationMiddleware struct {
	wrapped http.Handler
	session *session
}

// Returns a function which wraps around the router to provide authentification
func (s *session) genAuthHandler(wrapped http.Handler) authenticationMiddleware {
	return authenticationMiddleware{wrapped: wrapped, session: s}
}

// The actual logic behind it
// WARNING / WIP / TODO / TOKENTRANSITION: The token used by the client SHOULD NOT be the goil login cookie but should be independent.
func (h authenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check the header presence
	field := r.Header.Get("Authorization")
	if field == "" {
		err := NewError(AuthErrorAuthorizationHeaderAbsent, "", "").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	subFields := strings.Split(field, " ")

	// Check the header validity
	if len(subFields) != 2 {
		err := NewError(AuthErrorAuthorizationHeaderElementsLenNotTwo, "", "").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	} else if subFields[0] != "Bearer" { // Check if the authentification scheme is set to "Bearer"
		err := NewError(AuthErrorAuthorizationHeaderBearerTagAbsent, "", "").JSONWrite(w)
		if err != nil {
			panic(err)
		}
		return
	}

	// TODO / TOKENTRANSITION: Retrieve the cookie

	// Get a goil.Session using the cookie value given
	// TODO: Allow customization to the client options
	goilSession := goil.CreateSessionByCookieValue(subFields[1], &http.Client{})

	// Add to context
	ctx := context.WithValue(r.Context(), "session", goilSession)

	// If OK, then send it to the wrapped handler
	h.wrapped.ServeHTTP(w, r.WithContext(ctx))
}
