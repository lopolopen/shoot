package middleware

import (
	"log"
	"net/http"
	"time"
)

// RetryMiddleware returns a middleware that retries failed HTTP requests.
func RetryMiddleware(maxRetries int, delay time.Duration) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(req *http.Request) (*http.Response, error) {
			var resp *http.Response
			var err error

			// Attempt the request up to maxRetries + 1 times
			for attempt := 0; attempt <= maxRetries; attempt++ {
				if attempt > 0 {
					log.Printf("[Retry] Attempt %d for %s %s", attempt, req.Method, req.URL.String())
					time.Sleep(delay)
				}

				resp, err = next.RoundTrip(req)
				if err == nil && resp.StatusCode < 500 {
					return resp, nil
				}
			}

			return resp, err
		})
	}
}
