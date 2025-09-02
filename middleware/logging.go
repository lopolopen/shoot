package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is an HTTP client middleware that logs request details and response status.
func LoggingMiddleware(next http.RoundTripper) http.RoundTripper {
	return RoundTripper(func(req *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := next.RoundTrip(req)
		if err != nil {
			log.Printf("[HTTP] %s %s failed: %v (%v)", req.Method, req.URL.String(), err, time.Since(start))
			return nil, err
		}
		log.Printf("[HTTP] %s %s -> %d (%v)", req.Method, req.URL.String(), resp.StatusCode, time.Since(start))
		return resp, nil
	})
}
