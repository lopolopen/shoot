package main

import (
	"context"
	"fmt"
	"restclientsample/myclient"

	"github.com/lopolopen/shoot"
)

func main() {
	myC := shoot.NewRest[myclient.Client](
		shoot.BaseURLOfRestConf("http://localhost:8080"),
		shoot.TimeoutOfRestConf(1000),
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
