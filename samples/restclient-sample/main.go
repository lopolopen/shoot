package main

import (
	"context"
	"fmt"
	"net/http"
	"restclientsample/myclient"
	"time"

	"github.com/lopolopen/shoot"
	"github.com/lopolopen/shoot/middleware"
)

func WaitMiddeleware(next http.RoundTripper) http.RoundTripper {
	return middleware.RoundTripper(func(req *http.Request) (*http.Response, error) {
		time.Sleep(time.Second)
		return next.RoundTrip(req)
	})
}

func main() {
	myC := shoot.NewRest[myclient.Client](
		shoot.BaseURLOfRestConf("http://localhost:8080"),
		shoot.TimeoutOfRestConf(3000*time.Millisecond),
		shoot.EnableLoggingOfRestConf(true),
		shoot.Use(WaitMiddeleware),
	)

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
}
