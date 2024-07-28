package main

import (
	"net/http"
)

type BearerTransport struct {
	Token     string
	Transport http.RoundTripper
}

// RoundTrip implements http.RoundTripper interface.
func (b *BearerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	clonedRequest := r.Clone(r.Context())
	clonedRequest.Header.Set("Authorization", "Bearer "+b.Token)
	clonedRequest.Header.Set("Accept", "application/json")

	if b.Transport == nil {
		b.Transport = http.DefaultTransport
	}

	return b.Transport.RoundTrip(clonedRequest)
}
