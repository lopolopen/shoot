package main

import (
	"context"
	"fmt"
	"restclient_sample/gitlabclient"
)

func main() {
	ctx := context.Background()
	gitlabC := gitlabclient.New()
	user, err := gitlabC.GetUser(ctx, "12345")
	if err != nil {
		panic(err)
	}
	fmt.Printf("User ID: %s, Name: %s\n", user.ID, user.Name)
}
