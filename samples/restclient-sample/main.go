package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restclientsample/myclient"
	"time"

	"github.com/lopolopen/shoot"
)

func LoggingMiddleware(next http.RoundTripper) http.RoundTripper {
	return shoot.RoundTripper(func(req *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := next.RoundTrip(req)
		log.Printf("[HTTP] %s %s (%v)", req.Method, req.URL.String(), time.Since(start))
		return resp, err
	})
}

func main() {
	myC := shoot.NewRest[myclient.Client](
		shoot.BaseURLOfRestConf("http://localhost:8080"),
		shoot.TimeoutOfRestConf(3000*time.Millisecond),
		shoot.Use(LoggingMiddleware),
	).ConfigHTTPClient(func(c *http.Client) {
		//customize http.Client if needed
	})

	ctx := context.Background()
	err := myC.Set(ctx, myclient.NewKV("foo", "bar"))
	if err != nil {
		panic(err)
	}

	kv, err := myC.Get(ctx, "foo")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *kv)

	// xC := shoot.NewRest[xclient.Client](
	// 	shoot.BaseURLOfRestConf("http://localhost:8080"),
	// 	shoot.TimeoutOfRestConf(1000),
	// 	shoot.Use(LoggingMiddleware),
	// )
	// xC.NoComment(context.Background())
}
