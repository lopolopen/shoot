package middleware

import "net/http"

type Middleware func(next http.RoundTripper) http.RoundTripper

type RoundTripper func(req *http.Request) (*http.Response, error)

func (f RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
